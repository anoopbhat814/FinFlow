import type { Currency } from './account';

export type EntryDirection = 'DEBIT' | 'CREDIT';

export interface LedgerEntry {
  id: string;
  transactionId: string;
  accountId: string;
  direction: EntryDirection;
  amountMinor: number;
  currency: Currency;
  balanceAfterMinor: number;
  postedAt: string;
}

export interface LedgerTransaction {
  id: string;
  reference: string;
  description: string;
  entries: LedgerEntry[];
}

/** A ledger transaction is balanced when debits equal credits per currency. */
export function isBalanced(tx: LedgerTransaction): boolean {
  const totals = new Map<Currency, number>();
  for (const e of tx.entries) {
    const sign = e.direction === 'DEBIT' ? 1 : -1;
    totals.set(e.currency, (totals.get(e.currency) ?? 0) + sign * e.amountMinor);
  }
  return [...totals.values()].every((v) => v === 0);
}
