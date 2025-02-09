package main

import (
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.Reader
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	league, err := NewLeague(f.database)
	if err != nil {
		log.Println(err)
	}
	return league
}
