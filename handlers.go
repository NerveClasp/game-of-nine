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
