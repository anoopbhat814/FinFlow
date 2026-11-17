import { Money } from '@finflow/money';
import type { Currency } from '@finflow/types';

export function formatMinor(minor: number, currency: Currency, locale = 'en-US'): string {
  return Money.fromMinor(minor, currency).format(locale);
}

export function maskIban(iban?: string): string {
  if (!iban) return '—';
  return `${iban.slice(0, 4)} •••• ${iban.slice(-4)}`;
}

export function relativeTime(iso: string, now = Date.now()): string {
  const diffSec = Math.round((new Date(iso).getTime() - now) / 1000);
  const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' });
  const units: [Intl.RelativeTimeFormatUnit, number][] = [
    ['day', 86400],
    ['hour', 3600],
    ['minute', 60],
  ];
  for (const [unit, secs] of units) {
    if (Math.abs(diffSec) >= secs) return rtf.format(Math.round(diffSec / secs), unit);
  }
  return rtf.format(diffSec, 'second');
}
