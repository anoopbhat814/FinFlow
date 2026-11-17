import type { Account, Page } from '@finflow/types';
import { api } from '@/lib/api-client';
import { useAsync } from './useAsync';

export function useAccounts() {
  const state = useAsync(() => api<Page<Account>>('/accounts'));
  const totalsByCurrency = (state.data?.items ?? []).reduce<Record<string, number>>((acc, a) => {
    acc[a.currency] = (acc[a.currency] ?? 0) + a.balanceMinor;
    return acc;
  }, {});
  return { ...state, accounts: state.data?.items ?? [], totalsByCurrency };
}
