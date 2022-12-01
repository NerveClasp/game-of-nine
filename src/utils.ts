import { SORT_ORDER_KINDS, SORT_ORDER_VALUES } from './constants';
import type { CardType } from './types';

export const rand = (max = 1) => Math.floor((Math.random() * 100) % max);

export const sortHand = (cards: CardType[]): CardType[] => {
  const byKind: Array<Array<CardType>> = cards.reduce(
    (acc, card) => {
      const kindOrder = SORT_ORDER_KINDS[card.kind];
      acc[kindOrder].push(card);
      return acc;
    },
    [[], [], [], []]
  );

  return byKind
    .map((kind) =>
      kind.sort(
        (a, b) => SORT_ORDER_VALUES[a.value] - SORT_ORDER_VALUES[b.value]
      )
    )
    .flat();
};
