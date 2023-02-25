<script lang="ts">
	import type { BoardRowType } from './types';
	import Card from './Card.svelte';

	export let row: BoardRowType;
</script>

<div class="board-row">
	<span class="head">
		{#each row.head as card}
			<Card {card} />
		{/each}
	</span>
	<Card card={row.first} />
	{#each row.tail as card}
		<Card {card} />
	{/each}
</div>

<style>
	.board-row,
	.head {
		display: flex;
	}

	.head {
		flex-direction: row-reverse;
	}
</style>
