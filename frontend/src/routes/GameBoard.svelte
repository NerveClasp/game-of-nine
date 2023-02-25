<script lang="ts">
	import type { Game } from './types';
	import Card from '$lib/Card.svelte';

	export let game: Game;
</script>

<div class="game-board">
	{#each game.players as player}
		<div class="player">
			<span>{player.name}</span>
			<div class="player-cards">
				{#if !!player.cards}
					{#each player.cards as card}
						<Card {card} />
					{/each}
				{/if}
			</div>
		</div>
	{/each}
</div>

<style>
	.player-cards {
		display: flex;
		flex-wrap: wrap;
	}
</style>
