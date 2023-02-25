package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func containsType(message []byte, t string) bool {
	msgString := string(message)
	typeString := fmt.Sprintf("\"type\":\"%s\"", t)
	return strings.Contains(msgString, typeString)
}

func handleSocketMsg(message []byte, c *Client) {
	switch true {
	case containsType(message, "chat"):
		// do nothing
	case containsType(message, "createPlayer"):
		createPlayer(message)
	case containsType(message, "createGame"):
		createGame(message)
	case containsType(message, "startGame"):
		startGame(message, c)
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

func findGameIndex(m JoinOrLeaveGameMessage) int {
	for k, v := range games {
		if v.GameUid == m.GameUid {
			return k
		}
	}
	return -1
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

func startGame(message []byte, c *Client) {
	startMsg := JoinOrLeaveGameMessage{}
	err := json.Unmarshal(message, &startMsg)
	firstPlayerIsComputer := false
	computerPlayer := Player{}
	if err == nil {
		gameIdx := findGameIndex(startMsg)
		if gameIdx == -1 {
			return
		}
		shuffledCards := makeDeck()
		numberOfPlayers := len(games[gameIdx].Players)
		for kP := range games[gameIdx].Players {
			games[gameIdx].Players[kP].Cards = []Card{}
		}
		for k, v := range shuffledCards {
			playerIdx := k % numberOfPlayers
			games[gameIdx].Players[playerIdx].Cards = append(games[gameIdx].Players[playerIdx].Cards, v)
			if v.Kind == firstCard.Kind && v.Value == firstCard.Value {
				games[gameIdx].Players[playerIdx].YourTurn = true
				if games[gameIdx].Players[playerIdx].IsComputer {
					firstPlayerIsComputer = true
					computerPlayer = games[gameIdx].Players[playerIdx]
				}
			}
		}
		games[gameIdx].Started = true
		gameStarted := GameStarted{Type: "gameStarted", Game: games[gameIdx]}
		gameStartedMsg, err := json.Marshal(gameStarted)
		if err != nil {
			return
		}
		println(gameStartedMsg)
		c.hub.broadcast <- gameStartedMsg

		if firstPlayerIsComputer {
			playerMove := PlayerMoveMessage{Type: "playerMove", GameUid: games[gameIdx].GameUid, Player: computerPlayer}
			playerMoveMsg, err := json.Marshal(playerMove)
			if err != nil {
				return
			}
			c.hub.broadcast <- playerMoveMsg
		}
	}
}

// @FIX
func playerMove(message []byte) {

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
