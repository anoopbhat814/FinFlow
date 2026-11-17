'use client';
import { Card } from '@/components/Card';
import { BalanceCard } from '@/components/BalanceCard';
import { TransactionTable } from '@/components/TransactionTable';
import { useAccounts } from '@/hooks/useAccounts';
import { useTransfers } from '@/hooks/useTransfers';
import { formatMinor } from '@/lib/format';
import type { Currency } from '@finflow/types';

export default function DashboardPage() {
  const { accounts, totalsByCurrency, loading, error } = useAccounts();
  const transfers = useTransfers();

  if (loading) return <p>Loading your accounts…</p>;
  if (error) return <p role="alert">Could not load accounts: {error.message}</p>;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Overview</h1>
      <div className="flex flex-wrap gap-4 text-sm text-slate-600">
        {Object.entries(totalsByCurrency).map(([cur, total]) => (
          <span key={cur}>{cur} total: <strong>{formatMinor(total, cur as Currency)}</strong></span>
        ))}
      </div>
      <div className="grid gap-4 md:grid-cols-3">
        {accounts.map((a) => <BalanceCard key={a.id} account={a} />)}
      </div>
      <Card title="Recent transfers">
        <TransactionTable transfers={transfers.data?.items.slice(0, 8) ?? []} />
      </Card>
    </div>
  );
}
