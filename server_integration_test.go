package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	poker "webapp/v2"
)

func TestWinsRecordingAndRetrieval(t *testing.T) {
	database, cleanDatabase := poker.CreateTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := poker.NewFileSystemPlayerStore(database)
	poker.AssertNoError(t, err)

	server, err := poker.NewPlayerServer(store)
	poker.AssertNoError(t, err)

	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(t, player))

	poker.AssertStatusCode(t, response, http.StatusOK)

	poker.AssertResponseBody(t, response.Body.String(), "3")
}
