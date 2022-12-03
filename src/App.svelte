<script lang="ts">
  import Textfield, { HelperLine } from '@smui/textfield';
  import Switch from '@smui/switch';
  import Button from '@smui/button';
  import { STARTING_MONEY } from './constants';
  import type { Player } from './types';
  import Game from './lib/Game.svelte';

  let gameOn = false;
  let startingMoney = STARTING_MONEY;
  let allowNegative = false;

  $: moneyError = startingMoney < 1;
  $: canCreate =
    !moneyError &&
    players.length >= 3 &&
    players.every((p) => p.name.trim().length > 0);

  let players: Player[] = [];

  const handleAddPlayer = () => {
    players = [
      ...players,
      {
        name: '',
        cards: [],
        money: startingMoney,
        isComputer: false,
      },
    ];
  };

  const handleDeletePlayer = (idx: number) => {
    players.splice(idx, 1);
    players = [...players];
  };

  const handleCretateGame = () => {
    players = players.map(({ name, isComputer }) => ({
      name,
      isComputer: !!isComputer,
      cards: [],
      money: startingMoney,
    }));
    gameOn = true;
  };
</script>

<main style={`--players-count: ${players.length}`}>
  {#if gameOn}
    <Game bind:players bind:gameOn {allowNegative} />
  {:else}
    <section class="create-game">
      <h1>Create Game</h1>
      <Textfield
        bind:value={startingMoney}
        type="number"
        label="Starting money"
        invalid={moneyError}
      />
      <HelperLine>{moneyError ? 'Should be more than 1' : ' '}</HelperLine>
      <div>
        Allow players with 0 or less matches to play in a new draw:
        <Switch bind:checked={allowNegative} />
      </div>
      <h2>Players</h2>
      {#each players as player, idx}
        <div class="player">
          <Textfield
            kind="name"
            bind:value={player.name}
            label="Name"
            required
          />
          <span class="is-computer-label">Computer:</span>
          <Switch bind:checked={player.isComputer} />
          <Button on:click={() => handleDeletePlayer(idx)}>Delete</Button>
        </div>
      {/each}
      {#if players.length < 6}
        <Button class="material-icons" on:click={handleAddPlayer}>
          Add Player
        </Button>
      {/if}
      {#if players.length < 3}
        <div>Minimum 3 players are required</div>
      {/if}
      <div>
        <Button
          variant="raised"
          disabled={!canCreate}
          on:click={handleCretateGame}
        >
          Create Game
        </Button>
      </div>
    </section>
  {/if}
</main>

<style>
  :root {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    --players-count: 0;
  }

  * :global(.side-btn) {
    margin-left: 16px;
    margin-right: 16px;
  }

  * :global(.show-btn) {
    background-color: goldenrod;
    color: #000;
  }

  main {
    --rows-of-cards: calc(4 + var(--players-count, 0));
    display: flex;
    min-height: 100vh;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }

  .create-game {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
  }

  .is-computer-label {
    margin-left: 16px;
  }
</style>
