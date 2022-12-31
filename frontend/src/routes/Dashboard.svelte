<script lang="ts">
	import { onMount } from 'svelte';
	import { v4 as uuid } from 'uuid';
	import Button from '@smui/button/src/Button.svelte';
	import Textfield from '@smui/textfield';
	import socket, { chatMessages } from './socket';
	import type { Player, Game, IncomingMessage } from './types';

	export let player: Player;
	export let game: Game | null = null;

	let allMessages: IncomingMessage[] = [];
	let players: Player[] = [];
	let games: Game[] = [];
	let message: string = '';

	const sendMessage = () => {
		socket.sendMessage(message, player);
		message = '';
	};

	onMount(async () => {
		socket.subscribe((msg: IncomingMessage) => {
			if (!msg) return;
			if (msg?.type === 'createPlayer') players = [...players, msg.player];
			if (msg?.type === 'createGame') games = [...games, msg.game];

			allMessages = [...allMessages, msg];
		});

		players = await fetch('/get-players').then((r) => r.json());

		const existingGames = await fetch('/get-games')
			.then((r) => r.json())
			.catch((err) => console.log('getting games err:', err));
		if (Array.isArray(existingGames)) games = existingGames;
	});

	$: {
		console.log('allMessages:', allMessages);
	}

	const handleCreateGame = async () => {
		const gameUid = uuid();
		const newGame: Game = {
			players: [player],
			gameUid
		};
		socket.sendMessage({ type: 'createGame', game: newGame }, player);
		game = newGame;
	};

	const handleJoinGame = async (g: Game) => {
		if (!g.players.find((p) => p.playerUid === player.playerUid)) {
			socket.sendMessage({ type: 'joinGame', gameUid: g.gameUid }, player);
		}
		game = g;
	};
</script>

<h1>Dashboard</h1>
<Textfield bind:value={message} />
<Button on:click={sendMessage}>Send</Button>
<h2>Messages:</h2>
{#each $chatMessages as { player, message }}
	<div id={player.playerUid}>
		<span>{player.name}:</span><span>{message}</span>
	</div>
{/each}
<h2>Players:</h2>
{#each players as player}
	<div>{player.name}</div>
{/each}
<h2>Games:</h2>
{#each games as game}
	<div>{game.gameUid}<Button on:click={() => handleJoinGame(game)}>Join/Open</Button></div>
{/each}
<Button on:click={handleCreateGame}>Create Game</Button>
