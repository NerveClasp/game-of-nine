import type { CardValue, CardKind, CardType } from './types';

export const CARD_VALUES: CardValue[] = [
  '6',
  '7',
  '8',
  '9',
  '10',
  'J',
  'Q',
  'K',
  'A',
];
export const SORT_ORDER_VALUES = CARD_VALUES.reduce((acc, val, i) => {
  acc[val] = i;
  return acc;
}, {});

export const HEAD_ORDER: CardValue[] = ['8', '7', '6'];
export const TAIL_ORDER: CardValue[] = ['10', 'J', 'Q', 'K', 'A'];
export const CARD_KINDS: CardKind[] = ['❤', '♠', '♦', '♣'];
export const SORT_ORDER_KINDS = CARD_KINDS.reduce((acc, val, i) => {
  acc[val] = i;
  return acc;
}, {});

export const FIRST_CARD_VALUE: CardValue = '9';
export const FIRST_CARD: CardType = { kind: '❤', value: '9' };
export const STARTING_MONEY = 10;
export const MIN_ACTIVE_PLAYERS = 3;
