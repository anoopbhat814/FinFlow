import type { ButtonHTMLAttributes } from 'react';

type Variant = 'primary' | 'secondary' | 'danger';

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  loading?: boolean;
}

const styles: Record<Variant, string> = {
  primary: 'bg-indigo-600 text-white hover:bg-indigo-500',
  secondary: 'bg-white text-slate-900 border border-slate-300 hover:bg-slate-50',
  danger: 'bg-red-600 text-white hover:bg-red-500',
};

export function Button({ variant = 'primary', loading, disabled, children, className = '', ...rest }: Props) {
  return (
    <button
      {...rest}
      disabled={disabled || loading}
      aria-busy={loading}
      className={`rounded-md px-4 py-2 text-sm font-medium transition disabled:opacity-50 ${styles[variant]} ${className}`}
    >
      {loading ? 'Working…' : children}
    </button>
  );
}
