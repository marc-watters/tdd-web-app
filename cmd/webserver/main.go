package main

import (
	"log"
	"net/http"

	poker "webapp/v2"
)

const dbFilename = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000: %v", err)
	}
}
