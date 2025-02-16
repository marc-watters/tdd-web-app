package poker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	http.Handler

	game     Game
	store    PlayerStore
	template *template.Template
}

const htmlTemplatePath = "game.html"

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s: %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store
	p.game = game

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.gameHandler))
	router.Handle("/ws", http.HandlerFunc(p.wsHandler))

	p.Handler = router

	return p, nil
}

const JsonContentType = "application/json"

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", JsonContentType)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(p.store.GetLeague())
	if err != nil {
		log.Println(err)
	}
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("game.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("problem loading template: %s", err.Error()), http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, fmt.Sprintf("problem executing template: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (p *PlayerServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)
	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, err := strconv.Atoi(string(numberOfPlayersMsg))
	if err != nil {
		http.Error(w, fmt.Sprintf("problem parsing number of players input: %s", err.Error()), http.StatusInternalServerError)
	}
	p.game.Start(numberOfPlayers, ws)

	winnerMsg := ws.WaitForMsg()
	p.game.Finish(string(winnerMsg))
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

type playerServerWS struct {
	*websocket.Conn
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("problem upgrading connection to WebSocket: %v\n", err)
	}
	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading into websocket: %v\n", err)
	}
	return string(msg)
}

func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
