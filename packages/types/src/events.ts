import type { Transfer } from './transfer';
import type { Card } from './card';
import type { Invoice } from './invoice';

/** Domain events published on the notification stream. */
export type DomainEvent =
  | { type: 'transfer.completed'; transfer: Transfer }
  | { type: 'transfer.review_required'; transfer: Transfer }
  | { type: 'card.frozen'; card: Card }
  | { type: 'invoice.paid'; invoice: Invoice };

export function eventUserFacing(e: DomainEvent): string {
  switch (e.type) {
    case 'transfer.completed': return 'Transfer completed';
    case 'transfer.review_required': return 'Transfer under review';
    case 'card.frozen': return `Card ending ${e.card.last4} frozen`;
    case 'invoice.paid': return `Invoice ${e.invoice.number} paid`;
  }
}
