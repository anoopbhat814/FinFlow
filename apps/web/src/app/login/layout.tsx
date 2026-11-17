import type { Metadata } from 'next';
import type { ReactNode } from 'react';

export const metadata: Metadata = {
  title: 'Sign in · FinFlow',
};

export default function LoginLayout({ children }: { children: ReactNode }) {
  return <div className="mx-auto w-full max-w-md">{children}</div>;
}
