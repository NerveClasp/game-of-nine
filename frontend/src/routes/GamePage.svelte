<script lang="ts">
	import type { Game } from 'src/routes/types';

	export let game: Game;
</script>

<h1>Game</h1>
{#if game}
	<div>{JSON.stringify(game)}</div>
{/if}
