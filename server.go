package main

import (
	"fmt"
	"net/http"
	"strings"
)

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(player string) (score string) {
	switch player {
	case "Pepper":
		score = "20"
	case "Floyd":
		score = "10"
	default:
		score = ""
	}
	return
}
