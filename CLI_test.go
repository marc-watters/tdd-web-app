package poker_test

import (
	"strings"
	"testing"

	poker "webapp/v2"
)

var dummySpyBlindAlerter = &poker.SpyBlindAlerter{}

func TestCLI(t *testing.T) {
	t.Run("record Marc win from user input", func(t *testing.T) {
		in := strings.NewReader("Marc wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummySpyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Marc")
	})

	t.Run("record Chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in, dummySpyBlindAlerter)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Marc wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &poker.SpyBlindAlerter{}

		cli := poker.NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		if len(blindAlerter.Alerts) != 1 {
			t.Fatal("expected a blind alert to be scheduled")
		}
	})
}
