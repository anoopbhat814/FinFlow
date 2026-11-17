export type Currency = 'USD' | 'EUR' | 'GBP' | 'CAD' | 'JPY';
export type AccountStatus = 'PENDING' | 'ACTIVE' | 'FROZEN' | 'CLOSED';
export type AccountKind = 'CHECKING' | 'SAVINGS' | 'BUSINESS';

export interface Account {
  id: string;
  ownerId: string;
  name: string;
  kind: AccountKind;
  currency: Currency;
  status: AccountStatus;
  /** Settled balance in minor units (cents). */
  balanceMinor: number;
  /** Balance minus pending holds, in minor units. */
  availableMinor: number;
  iban?: string;
  createdAt: string;
}

export interface CreateAccountRequest {
  name: string;
  kind: AccountKind;
  currency: Currency;
}

export const isSpendable = (a: Account): boolean =>
  a.status === 'ACTIVE' && a.availableMinor > 0;
