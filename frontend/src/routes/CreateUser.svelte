<script lang="ts">
	import Textfield from '@smui/textfield';
	import Switch from '@smui/switch';
	import Button from '@smui/button';

	import type { NewPlayer } from './types';

	export let player: NewPlayer = { name: '', isComputer: false };
	export let handleAddPlayer: () => void;
</script>

<div class="player">
	<Textfield kind="name" bind:value={player.name} label="Name" required />
	<span class="is-computer-label">Computer:</span>
	<Switch bind:checked={player.isComputer} />
</div>
<Button on:click={handleAddPlayer} disabled={!player.name}>Create Player</Button>
