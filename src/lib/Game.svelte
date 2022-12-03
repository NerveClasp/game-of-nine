<script lang="ts">
  import { onMount } from 'svelte';
  import Button from '@smui/button';
  import Dialog, { Title, Content, Actions } from '@smui/dialog';
  import Card from './Card.svelte';
  import {
    CARD_KINDS,
    CARD_VALUES,
    FIRST_CARD,
    FIRST_CARD_VALUE,
    HEAD_ORDER,
    TAIL_ORDER,
    MIN_ACTIVE_PLAYERS,
  } from '../constants';

  import type { Player, CardType, BoardType } from '../types';
  import BoardRow from '../lib/BoardRow.svelte';
  import { rand, sortHand } from '../utils';
  import { play } from '../ai';

  export let players: Player[] = [];
  export let gameOn = false;
  export let allowNegative = false;

  let activePlayers = players;

  let showWinner = false;

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

  const onShuffle = () => {
    activePlayers.forEach((_p, i) => {
      activePlayers[i].cards = [];
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
      const playerIndex = i % activePlayers.length;
      const card = deck.pop();
      activePlayers[playerIndex].cards.push(card);
    }
    activePlayers = activePlayers.map((p) => {
      p.cards = sortHand(p.cards);
      return p;
    });
  };

  const reset = () => {
    toggleCards(false);
    resetBoard();
    if (!allowNegative) {
      activePlayers = activePlayers.filter(({ money }) => money > 0);
    }
    const gameEnded =
      activePlayers.reduce((acc, { money }) => {
        if (money > 0) acc += 1;
        return acc;
      }, 0) < (allowNegative ? 2 : MIN_ACTIVE_PLAYERS);
    if (gameEnded) {
      showWinner = true;
      return;
    }
    onShuffle();

    activePlayers.forEach((player, i) => {
      const firstCardIdx = player.cards.reduce((acc, { kind, value }, i) => {
        if (kind === FIRST_CARD.kind && value === FIRST_CARD.value) return i;
        return acc;
      }, -1);

      if (firstCardIdx > -1) {
        curPlayerIdx = i;
        activePlayers[i].cards[firstCardIdx].playable = true;
      }
    });
  };

  onMount(() => {
    reset();
  });

  const keyFromCard = ({ kind, value }: CardType): string => `${kind}${value}`;
  const nextPlayer = () => {
    toggleCards(false);
    curPlayerIdx = (curPlayerIdx + 1) % activePlayers.length;
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

    activePlayers[curPlayerIdx].cards = activePlayers[curPlayerIdx].cards.map(
      (card) => {
        const key = keyFromCard(card);
        if (!playableCards[key]) return card;
        return { ...card, playable: true };
      }
    );
    if (activePlayers[curPlayerIdx].isComputer) {
      makeAMove(play(activePlayers[curPlayerIdx].cards));
    }
  };

  const makeAMove = (card?: CardType) => {
    if (!card) {
      activePlayers[curPlayerIdx].cards = activePlayers[curPlayerIdx].cards.map(
        (c) => ({
          ...c,
          playable: false,
        })
      );
      activePlayers[curPlayerIdx].money -= 1;
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

    activePlayers[curPlayerIdx].cards = activePlayers[curPlayerIdx].cards
      .filter((c) => !(c.kind === kind && c.value === value))
      .map((c) => ({ ...c, playable: false }));

    if (activePlayers[curPlayerIdx].cards.length === 0) {
      activePlayers[curPlayerIdx].money += pot;
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

  const endGame = () => (gameOn = false);

  const getWinners = () => {
    const winners: Player[] = activePlayers.reduce((acc, player) => {
      if (!acc.length || acc[0].money < player.money) return [player];
      if (acc[0].money === player.money) acc.push(player);

      return acc as Player[];
    }, []);
    const label = `Winner${winners.length > 1 ? 's' : ''}:`;
    return `${label} ${winners.map(({ name }) => name).join(', ')}`;
  };
</script>

<Dialog open={showWinner} on:SMUIDialog:closed={endGame}>
  {#if showWinner}
    <Title>Game finished!</Title>
    <Content>{getWinners()}</Content>
    <Actions><Button on:click={endGame}>End Game</Button></Actions>
  {/if}
</Dialog>

<section class="heading">
  <h2 class="pot">Pot: {pot}</h2>
  <Button on:click={endGame}>End Game</Button>
</section>
<section class="board">
  {#each board as row}
    <BoardRow {row} />
  {/each}
</section>
<h3>activePlayers:</h3>
<section class="players">
  {#each activePlayers as { name, money, cards }, idx}
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

<style>
  :root {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    --players-count: 0;
  }

  * :global(.side-btn) {
    margin-left: 16px;
    margin-right: 16px;
    min-width: 12ch;
  }

  * :global(.show-btn) {
    background-color: goldenrod;
    color: #000;
  }

  h2,
  h3 {
    text-align: center;
    margin-top: 8px;
    margin-bottom: 8px;
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
