package poker_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	poker "webapp/v2"
)

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, &store)

	t.Run("return Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest(t, "Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatusCode(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("return Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest(t, "Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatusCode(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("return 404 on unknown players", func(t *testing.T) {
		request := newGetScoreRequest(t, "HeWhoShallNotBeNamed")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatusCode(t, response, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := mustMakePlayerServer(t, &store)

	t.Run("returns accepted on POST", func(t *testing.T) {
		player := "Pepper"

		request := newPostWinRequest(t, player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatusCode(t, response, http.StatusAccepted)
		poker.AssertPlayerWin(t, &store, player)
	})
}

func TestLeague(t *testing.T) {
	server := mustMakePlayerServer(t, &store)
	store := poker.StubPlayerStore{}

	t.Run("returns 200 on /league", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/league", nil)
		poker.AssertNoError(t, err)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []poker.Player

		err = json.NewDecoder(response.Body).Decode(&got)
		poker.AssertNoError(t, err)

		poker.AssertStatusCode(t, response, http.StatusOK)
	})

	t.Run("returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []poker.Player{
			{"Cleo", 32},
			{"Marc", 20},
			{"Tiest", 14},
		}

		server := mustMakePlayerServer(t, &store)
		store := poker.StubPlayerStore{nil, nil, wantedLeague}

		request := newLeagueRequest(t)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		poker.AssertStatusCode(t, response, http.StatusOK)
		poker.AssertContentType(t, response, poker.JsonContentType)
		poker.AssertLeague(t, got, wantedLeague)

		if response.Result().Header.Get("content-type") != "application/json" {
			t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
		}
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := poker.CreateTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)
	poker.AssertNoError(t, err)

	server := mustMakePlayerServer(t, store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(t, player))
		poker.AssertStatusCode(t, response, http.StatusOK)
		poker.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest(t))
		poker.AssertStatusCode(t, response, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []poker.Player{
			{"Pepper", 3},
		}
		poker.AssertLeague(t, got, want)
	})
}

func TestGame(t *testing.T) {
	t.Run("assert status code for /game endpoint", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{})

		request := newGameRequest(t)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertStatusCode(t, response, http.StatusOK)
	})

	t.Run("assert when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := StubPlayerStore{}
		winner := "Marc"
		server := httptest.NewServer(mustMakePlayerServer(t, &store))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
			t.Fatalf("could not send message over ws connection %v", err)
		}

		time.Sleep(10 * time.Millisecond)
		AssertPlayerWin(t, &store, winner)
	})
}

func newGetScoreRequest(t testing.TB, name string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	poker.AssertNoError(t, err)
	return req
}

func newPostWinRequest(t testing.TB, name string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	poker.AssertNoError(t, err)
	return req
}

func newLeagueRequest(t testing.TB) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/league", nil)
	poker.AssertNoError(t, err)
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []poker.Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("unable to parse response from server %q into slice of PLayer, '%v'", body, err)
	}
	return
}

func newGameRequest(t testing.TB) *http.Request {
	request, err := http.NewRequest(http.MethodGet, "/game", nil)
	poker.AssertNoError(t, err)
	return request
}

func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s: %v", url, err)
	}
	return ws
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection: %v", err)
	}
}
