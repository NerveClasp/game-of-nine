package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

type Card struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Player struct {
	Name       string `json:"name"`
	IsComputer bool   `json:"isComputer"`
	ClientUid  string `json:"clientUid"`
	PlayerUid  string `json:"playerUid"`
}

type Game struct {
	GameUid string   `json:"gameUid"`
	Players []Player `json:"players"`
}

type GetGame struct {
	GameUid string `json:"gameUid"`
}

type NewGameMessage struct {
	Type    string   `json:"type"`
	GameUid string   `json:"gameUid"`
	Players []Player `json:"players"`
}

var (
	players  []Player
	games    []Game
	upgrader = websocket.Upgrader{} // use default options
)

var (
	cardKinds    [4]string = [4]string{"❤", "♦", "♣", "♠"}
	cardValues   [9]string = [9]string{"6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	firstCard    Card      = Card{Kind: cardKinds[0], Value: cardValues[3]} // 9 of Hearts
	MIDDLE_VALUE string    = "9"
	HEAD_ORDER   [3]string = [3]string{"8", "7", "6"}
	TAIL_ORDER   [5]string = [5]string{"10", "J", "Q", "K", "A"}
)

func socket(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer connection.Close()
	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		msgString := string(message)

		log.Printf("recv: %s", message)
		isCreateGame := strings.Contains(msgString, "\"type\": \"create-game\"")
		if isCreateGame {
			log.Println("creating a game")
			newGame := NewGameMessage{}
			err := json.NewDecoder(r.Body).Decode(&newGame)
			if err != nil {
				msg := []byte("{ \"error\": \"could not parse message\" }")
				connection.WriteMessage(mt, msg)
				break
			}
			games = append(games, Game{GameUid: newGame.GameUid, Players: newGame.Players})
			msg := []byte("{ \"success\": true }")
			connection.WriteMessage(mt, msg)
		}
		err = connection.WriteMessage(mt, message)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func getGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS
	var getGame GetGame
	err := json.NewDecoder(r.Body).Decode(&getGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gameIdx := slices.IndexFunc(games, func(g Game) bool { return g.GameUid == getGame.GameUid })
	if gameIdx == -1 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games[gameIdx])
}

func addPlayer(w http.ResponseWriter, r *http.Request) {
	player := Player{}
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	players = append(players, player)
	w.WriteHeader(http.StatusOK)
}

func createGame(w http.ResponseWriter, r *http.Request) {
	game := Game{}
	err := json.NewDecoder(r.Body).Decode(&game)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	games = append(games, game)
	w.WriteHeader(http.StatusOK)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func lol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS
	w.WriteHeader(http.StatusOK)
	test := []string{}
	test = append(test, "Hello")
	test = append(test, "World")
	json.NewEncoder(w).Encode(test)
}

func main() {
	// Windows may be missing this
	mime.AddExtensionType(".js", "application/javascript")

	http.Handle("/lol", http.HandlerFunc(lol))
	http.Handle("/add-player", http.HandlerFunc(addPlayer))
	http.Handle("/get-players", http.HandlerFunc(getPlayers))
	http.Handle("/get-games", http.HandlerFunc(getGames))
	http.Handle("/get-game", http.HandlerFunc(getGame))
	http.Handle("/create-game", http.HandlerFunc(createGame))
	http.HandleFunc("/ws", socket)
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":7331", nil))
}
