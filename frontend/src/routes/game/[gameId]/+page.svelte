<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	import { page } from '$app/stores';
	import type { Game } from 'src/routes/types';

	let game: Game;

	const gameUid = $page.params.gameId;
	onMount(async () => {
		console.log({ gameUid });
		if (!gameUid) return goto('/');
		const resp = await fetch('/get-game', {
			method: 'POST',
			body: JSON.stringify({ gameUid })
		});
		if (resp.status !== 200) return goto('/');
		game = await resp.json();
	});
</script>

<h1>Game</h1>
{#if game}
	<div>{JSON.stringify(game?.players ?? [])}</div>
{/if}
