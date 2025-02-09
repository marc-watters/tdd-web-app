package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league   League
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	_, err := database.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("unable to seek database: %v\n", err)
	}
	league, _ := NewLeague(database)
	return &FileSystemPlayerStore{
		database: database,
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

	_, err := f.database.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("unable to seek database: %v\n", err)
	}
	err = json.NewEncoder(f.database).Encode(f.league)
	if err != nil {
		log.Printf("unable to write to database: %v", err)
	}
}
