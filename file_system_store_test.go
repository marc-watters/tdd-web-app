package poker_test

import (
	"testing"

	poker "webapp/v2"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Marc", "Wins": 20}
		]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		got := store.GetLeague()

		want := []poker.Player{
			{"Marc", 20},
			{"Cleo", 10},
		}

		poker.AssertLeague(t, got, want)
		// read again
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
				{"Name": "Cleo", "Wins": 10},
				{"Name": "Marc", "Wins": 20}
			]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		got := store.GetPlayerScore("Marc")
		want := 20

		poker.AssertScoreEquals(t, got, want)
	})

	t.Run("store wins for existing player", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Marc", "Wins": 20}
		]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		store.RecordWin("Marc")

		got := store.GetPlayerScore("Marc")
		want := 21

		poker.AssertScoreEquals(t, got, want)
	})

	t.Run("store wins for new player", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Marc", "Wins": 20}
		]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		poker.AssertScoreEquals(t, got, want)
	})

	t.Run("works with empty file", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := poker.CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Marc", "Wins": 20}
		]`)
		defer cleanDatabase()

		store, err := poker.NewFileSystemPlayerStore(database)
		poker.AssertNoError(t, err)

		got := store.GetLeague()

		want := poker.League{
			{"Marc", 20},
			{"Cleo", 10},
		}

		poker.AssertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		poker.AssertLeague(t, got, want)
	})
}
