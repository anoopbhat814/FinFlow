import { invoiceTotalMinor, isOverdue, type Invoice } from '@finflow/types';
import { formatMinor } from '@/lib/format';

export function InvoiceList({ invoices }: { invoices: Invoice[] }) {
  return (
    <ul className="divide-y divide-slate-100">
      {invoices.map((inv) => {
        const overdue = isOverdue(inv);
        return (
          <li key={inv.id} className="flex items-center justify-between py-3">
            <div>
              <p className="font-medium">{inv.number} · {inv.customerName}</p>
              <p className="text-xs text-slate-500">Due {new Date(inv.dueAt).toLocaleDateString()}</p>
            </div>
            <div className="text-right">
              <p className="tabular-nums">{formatMinor(invoiceTotalMinor(inv), inv.currency)}</p>
              <p className={`text-xs ${overdue ? 'text-red-600' : 'text-slate-500'}`}>
                {overdue ? 'overdue' : inv.status.toLowerCase()}
              </p>
            </div>
          </li>
        );
      })}
    </ul>
  );
}
