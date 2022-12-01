import type { CardType } from './types';
import { rand } from './utils';

export const play = (cards: CardType[]) => {
  const playable = cards.filter(({ playable }) => playable);

  if (!playable.length) return null;
  if (playable.length === 1) return playable[0];

  const randomCardIndex = rand(playable.length);

  return playable[randomCardIndex];
};
