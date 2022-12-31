<script lang="ts">
	import Button from '@smui/button';
	import { onMount } from 'svelte';
	import { v4 as uuid } from 'uuid';
	import socket, { clientUid } from './socket';

	import type { Game, IncomingMessage, Player } from './types';

	export let game: Game;

	onMount(() => {
		socket.subscribe((message: IncomingMessage) => {
			if (!message) return;
			if (message.gameUid !== game.gameUid) return;
			if (message.type === 'joinGame') game.players = [...game.players, message.player];
		});
	});

	const addComputer = async () => {
		const newComputer: Player = {
			name: 'PC',
			clientUid: $clientUid,
			playerUid: uuid(),
			isComputer: true
		};

		const resp = await fetch('/add-player', {
			method: 'POST', // *GET, POST, PUT, DELETE, etc.
			headers: {
				'Content-Type': 'application/json'
				// 'Content-Type': 'application/x-www-form-urlencoded',
			},
			body: JSON.stringify(newComputer)
		});
		if (resp.status !== 200) return; // @TODO handle error
		socket.sendMessage({ type: 'joinGame', gameUid: game.gameUid }, newComputer);
	};
</script>

<h2>Players:</h2>
{#each game.players as { name, playerUid }}
	<div id={playerUid}>{name}</div>
{/each}
<Button on:click={addComputer}>Add Computer player</Button>
