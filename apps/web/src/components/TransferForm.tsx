'use client';
import { useState, type FormEvent } from 'react';
import { Money } from '@finflow/money';
import type { Account, CreateTransferRequest, Transfer } from '@finflow/types';
import { validateTransfer, type FieldErrors } from '@/lib/validators';
import { useCreateTransfer } from '@/hooks/useTransfers';
import { Button } from './Button';

interface Props {
  accounts: Account[];
  onCreated?: (t: Transfer) => void;
}

export function TransferForm({ accounts, onCreated }: Props) {
  const [sourceId, setSourceId] = useState(accounts[0]?.id ?? '');
  const [destId, setDestId] = useState('');
  const [amount, setAmount] = useState('');
  const [memo, setMemo] = useState('');
  const [errors, setErrors] = useState<FieldErrors<CreateTransferRequest>>({});
  const [serverError, setServerError] = useState<string | null>(null);
  const { submit, submitting } = useCreateTransfer();

  const source = accounts.find((a) => a.id === sourceId);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!source) return;
    let amountMinor = 0;
    try {
      amountMinor = Money.parse(amount, source.currency).minor;
    } catch {
      setErrors({ amountMinor: 'Enter a valid amount' });
      return;
    }
    const req = { sourceAccountId: sourceId, destinationAccountId: destId, amountMinor, currency: source.currency, memo: memo || undefined };
    const found = validateTransfer(req, source.availableMinor);
    setErrors(found);
    if (Object.keys(found).length) return;
    try {
      onCreated?.(await submit(req));
      setAmount('');
      setMemo('');
      setServerError(null);
    } catch (err) {
      setServerError(err instanceof Error ? err.message : 'Transfer failed');
    }
  }

  return (
    <form onSubmit={onSubmit} className="space-y-3">
      <select value={sourceId} onChange={(e) => setSourceId(e.target.value)} className="w-full rounded border p-2">
        {accounts.map((a) => <option key={a.id} value={a.id}>{a.name} ({a.currency})</option>)}
      </select>
      <input value={destId} onChange={(e) => setDestId(e.target.value)} placeholder="Destination account ID" className="w-full rounded border p-2" />
      {errors.destinationAccountId && <p className="text-xs text-red-600">{errors.destinationAccountId}</p>}
      <input value={amount} onChange={(e) => setAmount(e.target.value)} inputMode="decimal" placeholder="0.00" className="w-full rounded border p-2" />
      {errors.amountMinor && <p className="text-xs text-red-600">{errors.amountMinor}</p>}
      <input value={memo} onChange={(e) => setMemo(e.target.value)} placeholder="Memo (optional)" className="w-full rounded border p-2" />
      {serverError && <p role="alert" className="text-sm text-red-600">{serverError}</p>}
      <Button type="submit" loading={submitting}>Send money</Button>
    </form>
  );
}
