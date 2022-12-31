export type NewPlayer = {
	name: string;
	isComputer: boolean;
};

export type Player = NewPlayer & {
	clientUid: string;
	playerUid: string;
	cards?: CardType[];
	cardsNumber?: number;
};

export type CardKind = '❤' | '♦' | '♣' | '♠';
export type CardValue = '6' | '7' | '8' | '9' | '10' | 'J' | 'Q' | 'K' | 'A';
export type CardType = {
	kind: CardKind;
	value: CardValue;
	hidden?: boolean;
	inactive?: boolean;
	playable?: boolean;
};

export type BoardRowType = {
	kind: CardKind;
	first: CardType;
	head: CardType[];
	tail: CardType[];
};

export type BoardType = BoardRowType[];

export type Game = {
	gameUid: string;
	players: Array<Player>;
	started?: boolean;
};

// @TODO: find out how to remove all `gameUid?:`
type PlayerMessage = { type: 'createPlayer'; gameUid?: string };
type MoveMessage = { type: 'move'; card?: CardType; gameUid: string };
type LobbyMessage = { type: 'joinLobby' | 'leaveLobby'; gameUid?: string };
type CreateGameMessage = { type: 'createGame'; game: Game; gameUid?: string };
type GameMessage = { type: 'joinGame' | 'leaveGame' | 'watchGame'; gameUid: string };
export type ChatMessage = { type: 'chat'; message: string; player: Player; gameUid?: string };

export type Message =
	| PlayerMessage
	| LobbyMessage
	| CreateGame
	| GameMessage
	| MoveMessage
	| ChatMessage;
export type OutgoingMessage = Message & { player: Player };
export type IncomingMessage = (Message & { player: Player }) | null;
