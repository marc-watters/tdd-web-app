package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type FileSystemPlayerStore struct {
	database io.Writer
	league   League
}

func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("unable to seek database: %v\n", err)
	}
	league, _ := NewLeague(file)
	return &FileSystemPlayerStore{
		database: &tape{file},
		league:   league,
	}
}

func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.league.Find(name)
	if player != nil {
		return player.Wins
	}
	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	err := json.NewEncoder(f.database).Encode(f.league)
	if err != nil {
		log.Printf("unable to write to database: %v", err)
	}
}
