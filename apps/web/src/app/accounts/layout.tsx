import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Accounts · FinFlow',
};

export default function AccountsLayout({ children }: { children: ReactNode }) {
  return <div className="mx-auto w-full max-w-5xl">{children}</div>;
}
