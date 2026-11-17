'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const LINKS = [
  { href: '/dashboard', label: 'Overview' },
  { href: '/accounts', label: 'Accounts' },
  { href: '/transfers', label: 'Transfers' },
  { href: '/invoices', label: 'Invoices' },
  { href: '/cards', label: 'Cards' },
];

export function Sidebar() {
  const pathname = usePathname();
  return (
    <nav aria-label="Primary" className="w-56 shrink-0 border-r border-slate-200 p-4">
      <div className="mb-6 text-lg font-bold text-indigo-700">FinFlow</div>
      <ul className="space-y-1">
        {LINKS.map((l) => {
          const active = pathname?.startsWith(l.href);
          return (
            <li key={l.href}>
              <Link
                href={l.href}
                aria-current={active ? 'page' : undefined}
                className={`block rounded-md px-3 py-2 text-sm ${active ? 'bg-indigo-50 font-medium text-indigo-700' : 'text-slate-600 hover:bg-slate-50'}`}
              >
                {l.label}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
