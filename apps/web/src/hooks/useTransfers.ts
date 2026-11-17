import { useCallback, useState } from 'react';
import type { CreateTransferRequest, Page, Transfer } from '@finflow/types';
import { api, newIdempotencyKey } from '@/lib/api-client';
import { useAsync } from './useAsync';

export function useTransfers(accountId?: string) {
  const q = accountId ? `?accountId=${encodeURIComponent(accountId)}` : '';
  return useAsync(() => api<Page<Transfer>>(`/transfers${q}`), [accountId]);
}

export function useCreateTransfer() {
  const [submitting, setSubmitting] = useState(false);
  const [key, setKey] = useState(newIdempotencyKey);

  const submit = useCallback(
    async (req: CreateTransferRequest) => {
      setSubmitting(true);
      try {
        const t = await api<Transfer>('/transfers', { method: 'POST', body: req, idempotencyKey: key });
        setKey(newIdempotencyKey());
        return t;
      } finally {
        setSubmitting(false);
      }
    },
    [key],
  );

  return { submit, submitting };
}
