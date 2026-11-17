import type { Currency } from './account';

export type CardType = 'VIRTUAL' | 'PHYSICAL';
export type CardStatus = 'INACTIVE' | 'ACTIVE' | 'FROZEN' | 'CANCELLED';

export interface Card {
  id: string;
  accountId: string;
  type: CardType;
  status: CardStatus;
  brand: 'VISA' | 'MASTERCARD';
  last4: string;
  expMonth: number;
  expYear: number;
  currency: Currency;
  monthlyLimitMinor: number;
  spentThisMonthMinor: number;
}

export interface UpdateCardLimitRequest {
  monthlyLimitMinor: number;
}

export function remainingLimit(card: Card): number {
  return Math.max(0, card.monthlyLimitMinor - card.spentThisMonthMinor);
}
