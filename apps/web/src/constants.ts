import type { Currency } from '@finflow/types';

export const SUPPORTED_CURRENCIES: readonly Currency[] = ['USD', 'EUR', 'GBP', 'CAD', 'JPY'];

export const NAV_ITEMS = [
  { href: '/dashboard', label: 'Overview' },
  { href: '/accounts', label: 'Accounts' },
  { href: '/transfers', label: 'Transfers' },
  { href: '/cards', label: 'Cards' },
  { href: '/invoices', label: 'Invoices' },
] as const;

export const POLL_INTERVAL_MS = 5_000;
export const MAX_MEMO_LENGTH = 140;
export const SESSION_WARNING_SECONDS = 120;
