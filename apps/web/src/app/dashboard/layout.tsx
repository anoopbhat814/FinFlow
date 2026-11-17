import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Overview · FinFlow',
};

export default function DashboardLayout({ children }: { children: ReactNode }) {
  return <div className="mx-auto w-full max-w-6xl">{children}</div>;
}
