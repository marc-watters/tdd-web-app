package main

import (
	"io"
	"log"
	"strings"
)

type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	_, err := f.database.Seek(0, io.SeekStart)
	if err != nil {
		log.Println(err)
	}

	league, err := NewLeague(f.database)
	if err != nil {
		log.Println(err)
	}
	return league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	var wins int
	for _, player := range f.GetLeague() {
		if player.Name == name {
			wins = player.Wins
			break
		}
	}
	return wins
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	database := strings.NewReader(`[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Marc", "Wins": 21}
	]`)

	f.database = database
}
