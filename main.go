package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"time"

	"golang.org/x/exp/slices"
)

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

var deckOfCards [len(cardValues) * len(cardKinds)]Card

func makeDeck() [36]Card {
	cardIndex := 0
	var suffledDeck [len(cardValues) * len(cardKinds)]Card
	for _, kind := range cardKinds {
		for _, value := range cardValues {
			deckOfCards[cardIndex] = Card{Kind: kind, Value: value}
			cardIndex++
		}
	}

	rand.Seed(time.Now().UnixNano())
	randIndexes := rand.Perm(len(deckOfCards))
	fmt.Printf("%+v\n", randIndexes)
	for k := range deckOfCards {
		suffledDeck[k] = deckOfCards[randIndexes[k]]
	}
	return suffledDeck
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
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":7331", nil))
}
