import type { Transfer } from '@finflow/types';

const escape = (v: string | number | undefined): string => {
  const s = v === undefined ? '' : String(v);
  return /[",\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s;
};

export function transfersToCsv(transfers: Transfer[]): string {
  const header = ['id', 'createdAt', 'status', 'currency', 'amount', 'memo'];
  const rows = transfers.map((t) => [
    t.id,
    t.createdAt,
    t.status,
    t.currency,
    (t.amountMinor / 100).toFixed(2),
    t.memo,
  ].map(escape).join(','));
  return [header.join(','), ...rows].join('\n');
}

export function downloadCsv(filename: string, csv: string): void {
  const url = URL.createObjectURL(new Blob([csv], { type: 'text/csv' }));
  const a = Object.assign(document.createElement('a'), { href: url, download: filename });
  a.click();
  URL.revokeObjectURL(url);
}
