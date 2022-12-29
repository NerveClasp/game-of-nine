<script lang="ts">
	import { onMount } from 'svelte';
	import CreateUser from './CreateUser.svelte';
	import { v4 as uuid } from 'uuid';
	import type { NewPlayer, Player } from './types';
	import Dashboard from './Dashboard.svelte';

	let playerCreated = false;
	let loading = false; // @TODO handle loading
	let newPlayer: NewPlayer = {
		name: '',
		isComputer: false
	};
	let player: Player;
	let players: Player[] = [];

	onMount(async () => {
		if (typeof window !== 'undefined') {
			loading = true;
			const playerUid = localStorage.getItem('playerUid');
			console.log('playerUid', playerUid);
			if (!playerUid) return;
			players = await fetch('/get-players').then((r) => r.json());
			console.log('players', players);
			const [existing] = players.filter((p: Player) => p.playerUid === playerUid);
			console.log('existing', existing);
			if (existing) {
				newPlayer = existing;
				playerCreated = true;
			}
			console.log('stop loading');
			loading = false;
		}
	});

	const handleAddPlayer = async () => {
		const clientUid = uuid();
		const playerUid = uuid();
		const p: Player = {
			...newPlayer,
			clientUid,
			playerUid
		};
		const addUser = await fetch('/add-player', {
			method: 'POST', // *GET, POST, PUT, DELETE, etc.
			headers: {
				'Content-Type': 'application/json'
				// 'Content-Type': 'application/x-www-form-urlencoded',
			},
			body: JSON.stringify(p)
		});
		// @TODO: handle error
		if (addUser.status !== 200) return;
		playerCreated = true;
		localStorage.setItem('clientUid', clientUid);
		localStorage.setItem('playerUid', playerUid);
		player = p;
	};
</script>

<svelte:head>
	<title>Game of Nine</title>
	<meta name="description" content="The card game of Nine" />
</svelte:head>

<section>
	<h1>Game of Nine</h1>
	{#if playerCreated}
		<Dashboard bind:player />
	{:else}
		<CreateUser bind:player={newPlayer} {handleAddPlayer} />
	{/if}
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 0.6;
	}

	h1 {
		width: 100%;
	}
</style>
