import type { Metadata } from 'next';
import type { ReactNode } from 'react';
import { Sidebar } from '@/components/Sidebar';

export const metadata: Metadata = {
  title: 'FinFlow',
  description: 'Move money with confidence.',
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body className="flex min-h-screen bg-slate-50 text-slate-900">
        <Sidebar />
        <main className="flex-1 p-8">{children}</main>
      </body>
    </html>
  );
}
