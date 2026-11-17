import type { Currency } from './account';
import type { FraudDecision } from './fraud';

export type TransferStatus =
  | 'PENDING'
  | 'IN_REVIEW'
  | 'COMPLETED'
  | 'FAILED'
  | 'REVERSED';

export interface Transfer {
  id: string;
  idempotencyKey: string;
  sourceAccountId: string;
  destinationAccountId: string;
  amountMinor: number;
  currency: Currency;
  memo?: string;
  status: TransferStatus;
  fraud?: FraudDecision;
  createdAt: string;
  completedAt?: string;
}

export interface CreateTransferRequest {
  sourceAccountId: string;
  destinationAccountId: string;
  amountMinor: number;
  currency: Currency;
  memo?: string;
}

export const TERMINAL_STATUSES: readonly TransferStatus[] = ['COMPLETED', 'FAILED', 'REVERSED'];

export function isTerminal(status: TransferStatus): boolean {
  return TERMINAL_STATUSES.includes(status);
}
