package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWinsRecordingAndRetrieval(t *testing.T) {
	database, cleanDatabase := CreateTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	AssertNoError(t, err)

	server := NewPlayerServer(store)

	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(t, player))

	AssertStatusCode(t, response, http.StatusOK)

	AssertResponseBody(t, response.Body.String(), "3")
}
