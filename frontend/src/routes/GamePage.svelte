<script lang="ts">
	import { onMount } from 'svelte';
	import socket from './socket';
	import type { Game } from 'src/routes/types';
	import Button from '@smui/button';
	import GameLobby from './GameLobby.svelte';

	export let game: Game;
	export let handleLeaveGame: (game: Game) => void;

	onMount(async () => {
		try {
			const resp = await fetch(`/get-game?gameUid=${game.gameUid}`);
			game = await resp.json();
		} catch (err) {
			console.log('Error fetching game:', err);
		}
	});
</script>

<h1>Game</h1>
<Button on:click={() => handleLeaveGame(game)}>Leave Game</Button>
{#if !game.started}
	<GameLobby bind:game />
{:else}
	<div>{JSON.stringify(game)}</div>
{/if}
