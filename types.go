package main

type Card struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Player struct {
	Name        string `json:"name"`
	IsComputer  bool   `json:"isComputer"`
	ClientUid   string `json:"clientUid"`
	PlayerUid   string `json:"playerUid"`
	Cards       []Card `json:"cards,omitempty"`
	CardsNumber int32  `json:"cardsNumber,omitempty"`
	YourTurn    bool   `json:"yourTurn,omitempty"`
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

type GameStarted struct {
	Type string `json:"type"`
	Game Game   `json:"game"`
}

type NewGameMessage struct {
	Type string `json:"type"`
	Game Game   `json:"game"`
}

type PlayerMoveMessage struct {
	Type    string `json:"type"` // playerMove
	GameUid string `json:"gameUid"`
	Player  Player `json:"player"`
	Card    Card   `json:"card"`
}

type GetGame struct {
	GameUid string `json:"gameUid"`
}

type JoinOrLeaveGameMessage struct {
	Type    string `json:"type"`
	GameUid string `json:"gameUid"`
	Player  Player `json:"player"`
}
