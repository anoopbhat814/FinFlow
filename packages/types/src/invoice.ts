import type { Currency } from './account';

export type InvoiceStatus = 'DRAFT' | 'SENT' | 'PAID' | 'OVERDUE' | 'VOID';

export interface InvoiceLine {
  description: string;
  quantity: number;
  unitPriceMinor: number;
}

export interface Invoice {
  id: string;
  number: string;
  customerName: string;
  customerEmail: string;
  currency: Currency;
  lines: InvoiceLine[];
  status: InvoiceStatus;
  issuedAt: string;
  dueAt: string;
  paidAt?: string;
}

export function invoiceTotalMinor(invoice: Pick<Invoice, 'lines'>): number {
  return invoice.lines.reduce((sum, l) => sum + l.quantity * l.unitPriceMinor, 0);
}

export function isOverdue(invoice: Invoice, now: Date = new Date()): boolean {
  return invoice.status === 'SENT' && new Date(invoice.dueAt).getTime() < now.getTime();
}
