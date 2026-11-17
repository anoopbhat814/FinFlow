import type { Transfer } from '@finflow/types';
import { formatMinor, relativeTime } from '@/lib/format';
import { FraudBadge } from './FraudBadge';

export function TransactionTable({ transfers }: { transfers: Transfer[] }) {
  if (transfers.length === 0) {
    return <p className="text-sm text-slate-500">No transfers yet.</p>;
  }
  return (
    <table className="w-full text-left text-sm">
      <thead className="text-xs uppercase text-slate-500">
        <tr>
          <th className="py-2">When</th>
          <th>Memo</th>
          <th>Status</th>
          <th>Risk</th>
          <th className="text-right">Amount</th>
        </tr>
      </thead>
      <tbody className="divide-y divide-slate-100">
        {transfers.map((t) => (
          <tr key={t.id}>
            <td className="py-2">{relativeTime(t.createdAt)}</td>
            <td>{t.memo ?? '—'}</td>
            <td>{t.status.replace('_', ' ').toLowerCase()}</td>
            <td><FraudBadge decision={t.fraud} /></td>
            <td className="text-right tabular-nums">{formatMinor(t.amountMinor, t.currency)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
