<script lang="ts">
	import { onMount } from 'svelte';
	import type { Player } from './types';

	// export let player: Player;
	let messages: string[] = [];
	let players: Player[] = [];

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
				messages = [...messages, event.data];
			};
		}

		players = await fetch('/get-players').then((r) => r.json());
		console.log('players', players);
		sendMessage('Hello from FE');
	});
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
