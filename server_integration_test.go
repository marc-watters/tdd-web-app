package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWinsRecordingAndRetrieval(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	server := NewPlayerServer(store)

	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(t, player))

	assertStatusCode(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}
