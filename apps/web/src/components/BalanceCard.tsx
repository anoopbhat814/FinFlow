import type { Account } from '@finflow/types';
import { formatMinor, maskIban } from '@/lib/format';
import { Card } from './Card';

export function BalanceCard({ account }: { account: Account }) {
  const frozen = account.status === 'FROZEN';
  return (
    <Card title={account.name}>
      <p className="text-3xl font-semibold tabular-nums">
        {formatMinor(account.balanceMinor, account.currency)}
      </p>
      <p className="mt-1 text-sm text-slate-500">
        Available {formatMinor(account.availableMinor, account.currency)} · {maskIban(account.iban)}
      </p>
      {frozen && (
        <p role="alert" className="mt-3 rounded bg-amber-50 px-3 py-2 text-sm text-amber-800">
          This account is frozen. Contact support to restore access.
        </p>
      )}
    </Card>
  );
}
