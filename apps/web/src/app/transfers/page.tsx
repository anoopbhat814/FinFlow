'use client';
import { Card } from '@/components/Card';
import { TransferForm } from '@/components/TransferForm';
import { TransactionTable } from '@/components/TransactionTable';
import { useAccounts } from '@/hooks/useAccounts';
import { useTransfers } from '@/hooks/useTransfers';

export default function TransfersPage() {
  const { accounts } = useAccounts();
  const transfers = useTransfers();
  return (
    <div className="grid gap-6 lg:grid-cols-[360px_1fr]">
      <Card title="New transfer">
        {accounts.length > 0 && <TransferForm accounts={accounts} onCreated={transfers.reload} />}
      </Card>
      <Card title="History">
        <TransactionTable transfers={transfers.data?.items ?? []} />
      </Card>
    </div>
  );
}
