package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strings"

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

type PlayerMessage struct {
	Type   string `json:"type"`
	Player Player `json:"player"`
}

type Game struct {
	GameUid string   `json:"gameUid"`
	Players []Player `json:"players"`
	Started bool     `json:"started,omitempty"`
}

type NewGameMessage struct {
	Type string `json:"type"`
	Game Game   `json:"game"`
}

type GetGame struct {
	GameUid string `json:"gameUid"`
}

type JoinOrLeaveGameMessage struct {
	Type    string `json:"type"`
	GameUid string `json:"gameUid"`
	Player  Player `json:"player"`
}

var (
	players []Player
	games   []Game
)

var (
	cardKinds    [4]string = [4]string{"❤", "♦", "♣", "♠"}
	cardValues   [9]string = [9]string{"6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	firstCard    Card      = Card{Kind: cardKinds[0], Value: cardValues[3]} // 9 of Hearts
	HEAD_ORDER   [3]string = [3]string{"8", "7", "6"}
	MIDDLE_VALUE string    = "9"
	TAIL_ORDER   [5]string = [5]string{"10", "J", "Q", "K", "A"}
)

func getGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS
	urlQuery := r.URL.Query()
	if !urlQuery.Has("gameUid") {
		http.Error(w, "No gameUid received", http.StatusBadRequest)
		return
	}

	gameUid := urlQuery.Get("gameUid")
	gameIdx := slices.IndexFunc(games, func(g Game) bool { return g.GameUid == gameUid })
	if gameIdx == -1 {
		http.Error(w, "Game not found", http.StatusBadRequest)
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

func containsType(message []byte, t string) bool {
	msgString := string(message)
	typeString := fmt.Sprintf("\"type\":\"%s\"", t)
	return strings.Contains(msgString, typeString)
}

func handleSocketMsg(message []byte) {
	switch true {
	case containsType(message, "chat"):
		// do nothing
	case containsType(message, "createPlayer"):
		createPlayer(message)
	case containsType(message, "createGame"):
		createGame(message)
	case containsType(message, "joinGame"):
		joinGame(message)
	case containsType(message, "leaveGame"):
		leaveGame(message)

	default:
		println("UNHANDLED")
		println(message)
	}
}

func createPlayer(message []byte) {
	playerMessage := PlayerMessage{}
	err := json.Unmarshal(message, &playerMessage)
	// @TODO: figure out how to handle errors for these puppies
	if err == nil {
		players = append(players, playerMessage.Player)
		// I'll keep this here and stop googling for it :)
		// fmt.Printf("%+v\n", playerMessage)
	}
}

func createGame(message []byte) {
	gameMsg := NewGameMessage{}
	err := json.Unmarshal(message, &gameMsg)
	if err == nil {
		games = append(games, gameMsg.Game)
	}
}

func joinGame(message []byte) {
	joinMsg := JoinOrLeaveGameMessage{}
	err := json.Unmarshal(message, &joinMsg)
	if err == nil {
		for k, v := range games {
			if v.GameUid == joinMsg.GameUid {
				games[k].Players = append(games[k].Players, joinMsg.Player)
			}
		}
	}
}

func leaveGame(message []byte) {
	leaveMsg := JoinOrLeaveGameMessage{}
	err := json.Unmarshal(message, &leaveMsg)
	if err == nil {
		for k, v := range games {
			if v.GameUid == leaveMsg.GameUid {
				newPlayers := []Player{}
				for _, p := range games[k].Players {
					if p.PlayerUid != leaveMsg.Player.PlayerUid {
						newPlayers = append(newPlayers, p)
					}
					games[k].Players = newPlayers
				}
			}
		}
	}
}

func createGameApi(w http.ResponseWriter, r *http.Request) {
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
	humanPlayers := []Player{}
	for _, v := range players {
		if !v.IsComputer {
			humanPlayers = append(humanPlayers, v)
		}
	}
	json.NewEncoder(w).Encode(humanPlayers)
}

func main() {
	// Windows may be missing this
	mime.AddExtensionType(".js", "application/javascript")
	hub := newHub()
	go hub.run()

	http.Handle("/add-player", http.HandlerFunc(addPlayer))
	http.Handle("/get-players", http.HandlerFunc(getPlayers))
	http.Handle("/get-games", http.HandlerFunc(getGames))
	http.Handle("/get-game", http.HandlerFunc(getGame))
	http.Handle("/create-game", http.HandlerFunc(createGameApi))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":7331", nil))
}
