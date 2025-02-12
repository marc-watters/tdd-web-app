package main

import (
	"log"
	"net/http"
	"os"

	poker "webapp/v2"
)

const dbFilename = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("unable to open db file %s: %v", dbFilename, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store: %v", err)
	}
	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000: %v", err)
	}
}
