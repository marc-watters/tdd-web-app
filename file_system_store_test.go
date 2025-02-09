package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Marc", "Wins": 20}
		]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Marc", 20},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database := strings.NewReader(`[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Marc", "Wins": 20}
			]`)

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Marc")
		want := 20

		if got != want {
			t.Errorf("\ngot: \t%d\nwant:\t%d", got, want)
		}
	})
}
