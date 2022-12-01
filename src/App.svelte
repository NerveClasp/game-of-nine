<script lang="ts">
  import { onMount } from 'svelte';
  import Button from '@smui/button';
  import Card from './lib/Card.svelte';
  import {
    CARD_KINDS,
    CARD_VALUES,
    FIRST_CARD,
    FIRST_CARD_VALUE,
    HEAD_ORDER,
    TAIL_ORDER,
    STARTING_MONEY,
  } from './constants';
  import type { CardType, BoardType } from './types';
  import BoardRow from './lib/BoardRow.svelte';
  import { rand, sortHand } from './utils';
  import { play } from './ai';

  const makeDeck = () =>
    CARD_KINDS.map((kind) =>
      CARD_VALUES.map((value) => ({ kind, value }))
    ).flat();
  let deck: CardType[] = makeDeck();
  let curPlayerIdx = 0;
  let curVisible = false;
  let pot = 0;

  let board: BoardType;
  const resetBoard = () => {
    const inactive = true;
    board = CARD_KINDS.map((kind) => {
      return {
        kind,
        first: {
          kind,
          value: FIRST_CARD_VALUE,
          inactive,
          playable: kind === FIRST_CARD.kind,
        },
        head: HEAD_ORDER.map((value) => ({ kind, value, inactive })),
        tail: TAIL_ORDER.map((value) => ({ kind, value, inactive })),
      };
    });
  };
  resetBoard();

  type Player = {
    name: string;
    cards: CardType[];
    money: number;
    isComputer?: boolean;
  };
  let players: Player[] = [
    { name: 'Romka', cards: [], money: STARTING_MONEY },
    { name: 'Elsa', cards: [], money: STARTING_MONEY },
    { name: 'PC', isComputer: true, cards: [], money: STARTING_MONEY },
  ];

  const onShuffle = () => {
    players.forEach((_p, i) => {
      players[i].cards = [];
    });

    deck = makeDeck();

    for (let i = 0; i < deck.length; i++) {
      const temp = deck[i];
      const randomIndex = rand(deck.length);
      deck[i] = deck[randomIndex];
      deck[randomIndex] = temp;
    }

    const deckLength = deck.length;
    for (let i = 0; i < deckLength; i++) {
      const playerIndex = i % players.length;
      const card = deck.pop();
      players[playerIndex].cards.push(card);
    }
    players = players.map((p) => {
      p.cards = sortHand(p.cards);
      return p;
    });
  };

  const reset = () => {
    toggleCards(false);
    resetBoard();
    onShuffle();

    players.forEach((player, i) => {
      const firstCardIdx = player.cards.reduce((acc, { kind, value }, i) => {
        if (kind === FIRST_CARD.kind && value === FIRST_CARD.value) return i;
        return acc;
      }, -1);

      if (firstCardIdx > -1) {
        curPlayerIdx = i;
        players[i].cards[firstCardIdx].playable = true;
      }
    });
  };
  onMount(() => {
    reset();
  });

  const keyFromCard = ({ kind, value }: CardType): string => `${kind}${value}`;
  const nextPlayer = () => {
    toggleCards(false);
    curPlayerIdx = (curPlayerIdx + 1) % players.length;
    const playableCards = board.reduce((acc, row) => {
      if (row.first.playable) {
        acc[keyFromCard(row.first)] = true;
        return acc;
      }
      row.head.forEach((card) => {
        if (card.playable) acc[keyFromCard(card)] = true;
      });
      row.tail.forEach((card) => {
        if (card.playable) acc[keyFromCard(card)] = true;
      });
      return acc;
    }, {});

    players[curPlayerIdx].cards = players[curPlayerIdx].cards.map((card) => {
      const key = keyFromCard(card);
      if (!playableCards[key]) return card;
      return { ...card, playable: true };
    });
    if (players[curPlayerIdx].isComputer) {
      makeAMove(play(players[curPlayerIdx].cards));
    }
  };

  const makeAMove = (card?: CardType) => {
    if (!card) {
      players[curPlayerIdx].cards = players[curPlayerIdx].cards.map((c) => ({
        ...c,
        playable: false,
      }));
      players[curPlayerIdx].money -= 1;
      pot += 1;
      return nextPlayer();
    }

    const { value, kind } = card;
    const valueIsFirst = value === FIRST_CARD.value;
    const kindIsFirst = kind === FIRST_CARD.kind;
    const isFirstCardOfTheGame = valueIsFirst && kindIsFirst;
    const valueIsHead = !valueIsFirst && HEAD_ORDER.includes(value);

    board = board.map((row) => {
      if (row.kind !== kind) {
        if (!isFirstCardOfTheGame) return row;
        row.first.playable = true;
        return row;
      }

      if (valueIsFirst) {
        row.first.playable = false;
        row.first.inactive = false;
        row.head[0].playable = true;
        row.tail[0].playable = true;
      } else if (valueIsHead) {
        let found = false;
        for (let i = 0; i < row.head.length; i++) {
          if (found) {
            row.head[i].playable = true;
            break;
          }
          if (row.head[i].playable) {
            row.head[i].playable = false;
            row.head[i].inactive = false;
            found = true;
          }
        }
      } else {
        let found = false;
        for (let i = 0; i < row.tail.length; i++) {
          if (found) {
            row.tail[i].playable = true;
            break;
          }
          if (row.tail[i].playable) {
            row.tail[i].playable = false;
            row.tail[i].inactive = false;
            found = true;
          }
        }
      }
      return row;
    });
    players[curPlayerIdx].cards = players[curPlayerIdx].cards
      .filter((c) => !(c.kind === kind && c.value === value))
      .map((c) => ({ ...c, playable: false }));
    if (players[curPlayerIdx].cards.length === 0) {
      players[curPlayerIdx].money += pot;
      pot = 0;
      reset();
    } else {
      nextPlayer();
    }
  };

  const toggleCards = (visible: boolean) => {
    curVisible = !!visible;
  };

  const formatPlayerInfo = ({
    name,
    money,
  }: {
    name: string;
    money: number;
  }): string => `${name} | $${money}`;
</script>

<main style={`--players-count: ${players.length}`}>
  <h2 class="pot">Pot: {pot}</h2>
  <section class="board">
    {#each board as row}
      <BoardRow {row} />
    {/each}
  </section>
  <h3>Players:</h3>
  <section class="players">
    {#each players as { name, money, cards }, idx}
      <div class="cards-row">
        {#if curPlayerIdx === idx}
          <span class="current">
            {formatPlayerInfo({ name, money })}
          </span>
          <div class="player-cards">
            <Button
              class="side-btn show-btn"
              variant="raised"
              on:click={() => toggleCards(true)}>Show</Button
            >
            {#each cards as card}
              <Card
                {card}
                on:click={() => {
                  if (!curVisible) return toggleCards(true);
                  if (!card.playable) return;
                  makeAMove(card);
                }}
                hidden={!curVisible}
                clickable={card.playable || !curVisible}
              />
            {/each}
            <Button
              class="side-btn"
              variant="raised"
              color="secondary"
              on:click={() => makeAMove()}>Pay $</Button
            >
          </div>
        {:else}
          <span class="inactive">
            {formatPlayerInfo({ name, money })}
          </span>
          <div class="player-cards">
            {#each cards as card}
              <Card {card} hidden />
            {/each}
          </div>
        {/if}
      </div>
    {/each}
  </section>
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
  }

  h2,
  h3 {
    text-align: center;
  }

  .board {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .players {
    display: flex;
    flex-direction: column;
    text-align: center;
    align-items: center;
    gap: 4px;
  }

  .player-cards {
    display: flex;
    width: 100%;
    align-items: center;
  }

  .current {
    font-weight: 700;
  }

  .inactive {
    opacity: 0.6;
  }
</style>
