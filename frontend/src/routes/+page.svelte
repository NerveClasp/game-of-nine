<script lang="ts">
	import { onMount } from 'svelte';
	import CreateUser from './CreateUser.svelte';
	import { v4 as uuid } from 'uuid';
	import type { Game, IncomingMessage, NewPlayer, Player } from './types';
	import Dashboard from './Dashboard.svelte';
	import GamePage from './GamePage.svelte';
	import socket, { clientUid as clientUidStore } from './socket';

	let playerCreated = false;
	let loading = false; // @TODO handle loading
	let newPlayer: NewPlayer = {
		name: '',
		isComputer: false
	};
	let player: Player;
	let players: Player[] = [];
	let game: Game | null = null;

	onMount(async () => {
		if (typeof window !== 'undefined') {
			loading = true;
			const playerUid = localStorage.getItem('playerUid');
			if (!playerUid) return;

			players = await fetch('/get-players').then((r) => r.json());
			const [existing] = players?.filter((p: Player) => p.playerUid === playerUid) ?? [];
			if (existing) {
				player = existing;
				playerCreated = true;
			}
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
		socket.sendMessage({ type: 'createPlayer' }, p);
		// @TODO: handle error
		playerCreated = true;
		localStorage.setItem('clientUid', clientUid);
		localStorage.setItem('playerUid', playerUid);
		$clientUidStore = clientUid;
		player = p;
	};

	const handleLeaveGame = async (g: Game) => {
		socket.sendMessage({ type: 'leaveGame', gameUid: g.gameUid }, player);
		game = null;
	};
</script>

<svelte:head>
	<title>Game of Nine</title>
	<meta name="description" content="The card game of Nine" />
</svelte:head>

<section>
	<h1>Game of Nine</h1>
	{#if playerCreated}
		{#if !!game}
			<GamePage bind:game {handleLeaveGame} />
		{:else}
			<Dashboard bind:player bind:game />
		{/if}
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
