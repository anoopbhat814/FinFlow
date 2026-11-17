import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Invoices · FinFlow',
};

export default function InvoicesLayout({ children }: { children: ReactNode }) {
  return <div className="mx-auto w-full max-w-5xl">{children}</div>;
}
