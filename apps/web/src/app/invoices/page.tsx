'use client';
import type { Invoice, Page } from '@finflow/types';
import { Card } from '@/components/Card';
import { InvoiceList } from '@/components/InvoiceList';
import { api } from '@/lib/api-client';
import { useAsync } from '@/hooks/useAsync';

export default function InvoicesPage() {
  const { data, loading, error } = useAsync(() => api<Page<Invoice>>('/invoices'));
  return (
    <Card title="Invoices">
      {loading && <p>Loading…</p>}
      {error && <p role="alert">{error.message}</p>}
      {data && <InvoiceList invoices={data.items} />}
    </Card>
  );
}
