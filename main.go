package main

import (
	"log"
	"net/http"
)

type InMemoryPlayerStore struct{}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}
	log.Fatal(http.ListenAndServe(":5000", server))
}
