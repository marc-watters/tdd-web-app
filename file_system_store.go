package main

import (
	"io"
	"log"
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
