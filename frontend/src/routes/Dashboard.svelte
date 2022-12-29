<script lang="ts">
	import { onMount } from 'svelte';
	import { v4 as uuid } from 'uuid';
	import Button from '@smui/button/src/Button.svelte';
	import type { Player, Game } from './types';

	export let player: Player;
	export let game: Game | null = null;

	let messages: string[] = [];
	let players: Player[] = [];
	let games: Game[] = [];

	console.log('player', player);

	let socket: WebSocket | null = null;

	const sendMessage = (message: string) => {
		console.log('socket', socket);
		if (socket && socket.readyState <= 1) {
			socket.send(message);
		}
	};

	onMount(async () => {
		if (typeof window !== 'undefined') {
			socket = new WebSocket('ws://localhost:7331/ws');

			// Connection opened
			socket.onopen = (_event: Event) => {
				console.log("It's open");
			};

			// Listen for messages
			socket.onmessage = (event) => {
				console.log('event', event);
				messages = [...messages, event.data];
			};
		}

		players = await fetch('/get-players').then((r) => r.json());
		console.log('players', players);
		const existingGames = await fetch('/get-games')
			.then((r) => r.json())
			.catch((err) => console.log('getting games err:', err));
		if (Array.isArray(existingGames)) games = existingGames;
		console.log('games', games);
		sendMessage('Hello from FE');
	});

	const handleCreateGame = async () => {
		const gameUid = uuid();
		const newGame: Game = {
			players: [player],
			gameUid
		};
		const createGame = await fetch('/create-game', {
			method: 'POST', // *GET, POST, PUT, DELETE, etc.
			headers: {
				'Content-Type': 'application/json'
				// 'Content-Type': 'application/x-www-form-urlencoded',
			},
			body: JSON.stringify(newGame)
		});
		if (createGame.status !== 200) return; // @TODO handle error
		game = newGame;
	};
</script>

<h1>Dashboard</h1>
<h2>Messages:</h2>
{#each messages as message}
	<div>{message}</div>
{/each}
<h2>Players:</h2>
{#each players as player}
	<div>{player.name}</div>
{/each}
<h2>Games:</h2>
{#each games as game}
	<div>{game.gameUid}</div>
{/each}
<Button on:click={handleCreateGame}>Create Game</Button>
