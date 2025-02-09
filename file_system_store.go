package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.Reader
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
	var league []Player
	err := json.NewDecoder(f.database).Decode(&league)
	if err != nil {
		log.Println(err)
	}
	return league
}
