package poker

import (
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	return s.Scores[name]
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.WinCalls = append(s.WinCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.League
}

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	s.Alerts = append(s.Alerts, ScheduledAlert{duration, amount})
}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

type SpyGame struct {
	StartCalled     bool
	StartCalledWith int

	BlindAlert []byte

	FinishCalled     bool
	FinishCalledWith string
}

func (s *SpyGame) Start(numberOfPlayers int, out io.Writer) {
	s.StartCalled = true
	s.StartCalledWith = numberOfPlayers
	_, err := out.Write(s.BlindAlert)
	if err != nil {
		log.Printf("error writing to out: %v", err)
	}
}

func (s *SpyGame) Finish(winner string) {
	s.FinishCalled = true
	s.FinishCalledWith = winner
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.WinCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
	}

	if store.WinCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.WinCalls[0], winner)
	}
}

func AssertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot: \t%q\nwant:\t%q", got, want)
	}
}

func AssertStatusCode(t testing.TB, got *httptest.ResponseRecorder, want int) {
	t.Helper()
	if got.Code != want {
		t.Errorf("\ngot: \t%d\nwant:\t%d", got.Code, want)
	}
}

func AssertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: \t%v\nwant:\t%v", got, want)
	}
}

func AssertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func AssertScoreEquals(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot: \t%d\nwant:\t%d", got, want)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}
}

func CreateTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	_, err = tmpfile.Write([]byte(initialData))
	AssertNoError(t, err)

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
