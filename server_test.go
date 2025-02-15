package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("return Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest(t, "Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatusCode(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("return Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest(t, "Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatusCode(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("return 404 on unknown players", func(t *testing.T) {
		request := newGetScoreRequest(t, "HeWhoShallNotBeNamed")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatusCode(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("returns accepted on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(t, player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatusCode(t, response.Code, http.StatusAccepted)
		AssertPlayerWin(t, &store, player)
	})
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("returns 200 on /league", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/league", nil)
		AssertNoError(t, err)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player

		err = json.NewDecoder(response.Body).Decode(&got)
		AssertNoError(t, err)

		AssertStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Marc", 20},
			{"Tiest", 14},
		}

		store := &StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(store)

		request := newLeagueRequest(t)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		AssertStatusCode(t, response.Code, http.StatusOK)
		AssertContentType(t, response, jsonContentType)
		AssertLeague(t, got, wantedLeague)

		if response.Result().Header.Get("content-type") != "application/json" {
			t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := CreateTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	AssertNoError(t, err)

	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(t, player))
		AssertStatusCode(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest(t))
		AssertStatusCode(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		AssertLeague(t, got, want)
	})
}

func TestGame(t *testing.T) {
	server := NewPlayerServer(&StubPlayerStore{})

	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	AssertStatusCode(t, response.Code, http.StatusOK)
}

func newGetScoreRequest(t testing.TB, name string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	AssertNoError(t, err)
	return req
}

func newPostWinRequest(t testing.TB, name string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	AssertNoError(t, err)
	return req
}

func newLeagueRequest(t testing.TB) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/league", nil)
	AssertNoError(t, err)
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("unable to parse response from server %q into slice of PLayer, '%v'", body, err)
	}
	return
}
