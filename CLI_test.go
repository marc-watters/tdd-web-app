package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	poker "webapp/v2"
)

var (
	dummySpyBlindAlerter = &poker.SpyBlindAlerter{}
	dummyPlayerStore     = &poker.StubPlayerStore{}
	dummyStdOut          = &bytes.Buffer{}
)

func TestCLI(t *testing.T) {
	t.Run("start game with 3 players and finish game with 'Marc' as winner", func(t *testing.T) {
		game := &poker.SpyGame{}
		stdout := &bytes.Buffer{}

		in := userSends("3", "Marc wins")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Marc")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &poker.SpyGame{}

		in := userSends("8", "Cleo wins")
		cli := poker.NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("NaN\n")
		game := &poker.SpyGame{}

		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Errorf("game should not have started")
		}

		wantPrompt := poker.PlayerPrompt + poker.ErrPlayerInputMsg
		assertMessagesSentToUser(t, stdout, wantPrompt)
	})

	t.Run("it prints an error when the winner is declared incorrectly", func(t *testing.T) {
		game := &poker.SpyGame{}
		stdout := &bytes.Buffer{}

		in := userSends("8", "Lloyd is a killer")
		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertGameNotFinished(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.ErrWinnerInputMsg)
	})
}

func assertScheduledAlert(t testing.TB, got, want poker.ScheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func checkSchedulingCases(cases []poker.ScheduledAlert, t testing.TB, alerter *poker.SpyBlindAlerter) {
	t.Helper()
	for i, want := range cases {
		if len(alerter.Alerts) <= i {
			t.Fatalf("alert %d was not scheduled %v", i, alerter.Alerts)
		}

		got := alerter.Alerts[i]
		assertScheduledAlert(t, got, want)
	}
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	got := stdout.String()
	want := strings.Join(messages, "")

	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func assertGameStartedWith(t testing.TB, game *poker.SpyGame, numberOfPlayersWanted int) {
	t.Helper()
	if game.StartCalledWith != numberOfPlayersWanted {
		t.Errorf("wanted Start called with %d but got %d", numberOfPlayersWanted, game.StartCalledWith)
	}
}

func assertFinishCalledWith(t testing.TB, game *poker.SpyGame, winner string) {
	t.Helper()
	if game.FinishCalledWith != winner {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishCalledWith)
	}
}

func assertGameNotFinished(t testing.TB, game *poker.SpyGame) {
	t.Helper()
	if game.FinishCalled {
		t.Errorf("game should not have finished")
	}
}

func userSends(messages ...string) io.Reader {
	return strings.NewReader(strings.Join(messages, "\n"))
}
