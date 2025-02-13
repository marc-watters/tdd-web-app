package main

import (
	"fmt"
	"log"
	"os"
	poker "webapp/v2"
)

const dbFileName = "game.db.json"

var dummyBlindAlerter poker.BlindAlerter = nil

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {NAME} wins to record a win")
	poker.NewCLI(store, os.Stdin, dummyBlindAlerter).PlayPoker()
}
