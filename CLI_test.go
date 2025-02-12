package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	in := strings.NewReader("Marc wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	assertPlayerWin(t, playerStore, "Marc")
}
