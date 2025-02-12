package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {
	scanner := bufio.NewScanner(cli.in)
	scanner.Scan()
	cli.playerStore.RecordWin(extractWinner(scanner.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
