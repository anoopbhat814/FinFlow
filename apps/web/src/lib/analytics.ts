import type { Transfer, Currency } from '@finflow/types';

export interface DailyTotal { date: string; outflowMinor: number; count: number }

/** Groups completed transfers by UTC day for the spending chart. */
export function dailyOutflow(transfers: Transfer[], currency: Currency): DailyTotal[] {
  const byDay = new Map<string, DailyTotal>();
  for (const t of transfers) {
    if (t.status !== 'COMPLETED' || t.currency !== currency) continue;
    const date = t.createdAt.slice(0, 10);
    const cur = byDay.get(date) ?? { date, outflowMinor: 0, count: 0 };
    cur.outflowMinor += t.amountMinor;
    cur.count += 1;
    byDay.set(date, cur);
  }
  return [...byDay.values()].sort((a, b) => a.date.localeCompare(b.date));
}

export function movingAverage(values: number[], window: number): number[] {
  return values.map((_, i) => {
    const slice = values.slice(Math.max(0, i - window + 1), i + 1);
    return slice.reduce((a, b) => a + b, 0) / slice.length;
  });
}
