export type NewPlayer = {
	name: string;
	isComputer: boolean;
};

export type Player = NewPlayer & {
	clientUid: string;
	playerUid: string;
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
};
