'use client';
import { BalanceCard } from '@/components/BalanceCard';
import { useAccounts } from '@/hooks/useAccounts';

export default function AccountsPage() {
  const { accounts, loading, error } = useAccounts();
  if (loading) return <p>Loading…</p>;
  if (error) return <p role="alert">{error.message}</p>;
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Accounts</h1>
      <div className="grid gap-4 md:grid-cols-2">
        {accounts.map((a) => <BalanceCard key={a.id} account={a} />)}
      </div>
    </div>
  );
}
