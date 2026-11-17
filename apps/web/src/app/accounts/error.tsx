'use client';
import { Button } from '@/components/Button';

interface Props {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function AccountsError({ error, reset }: Props) {
  return (
    <div role="alert" className="rounded-xl border border-red-200 bg-red-50 p-6">
      <h2 className="text-base font-semibold text-red-800">We couldn't load your accounts</h2>
      <p className="mt-1 text-sm text-red-700">{error.message}</p>
      {error.digest && <p className="mt-1 text-xs text-red-500">Reference: {error.digest}</p>}
      <Button variant="secondary" className="mt-4" onClick={reset}>
        Try again
      </Button>
    </div>
  );
}
