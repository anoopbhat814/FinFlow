import type { ReactNode } from 'react';

interface CardProps {
  title?: string;
  action?: ReactNode;
  children: ReactNode;
}

export function Card({ title, action, children }: CardProps) {
  return (
    <section className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
      {(title || action) && (
        <header className="mb-4 flex items-center justify-between">
          {title && <h2 className="text-base font-semibold text-slate-900">{title}</h2>}
          {action}
        </header>
      )}
      {children}
    </section>
  );
}
