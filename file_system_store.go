package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
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
	league := f.GetLeague()

	for i, player := range league {
		if player.Name == name {
			league[i].Wins++
		}
	}

	_, err := f.database.Seek(0, io.SeekStart)
	if err != nil {
		log.Println(err)
	}

	err = json.NewEncoder(f.database).Encode(league)
	if err != nil {
		log.Println(err)
	}
}
