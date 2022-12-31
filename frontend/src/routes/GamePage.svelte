<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import type { Game } from 'src/routes/types';
	import Button from '@smui/button';
	import GameLobby from './GameLobby.svelte';

	export let game: Game;
	export let handleLeaveGame: (game: Game) => void;
	export let handleBackToLobby: () => void;

	onMount(async () => {
		if (typeof window !== 'undefined') {
			localStorage.setItem('gameUid', game.gameUid);
		}
		try {
			const resp = await fetch(`/get-game?gameUid=${game.gameUid}`);
			game = await resp.json();
		} catch (err) {
			console.log('Error fetching game:', err);
		}
	});

	onDestroy(() => {
		if (typeof window !== 'undefined') {
			localStorage.removeItem('gameUid');
		}
	});
</script>

<h1>Game</h1>
<Button on:click={() => handleLeaveGame(game)}>Leave Game</Button>
<Button on:click={handleBackToLobby}>Back to games list</Button>
{#if !game.started}
	<GameLobby bind:game />
{:else}
	<div>{JSON.stringify(game)}</div>
{/if}
