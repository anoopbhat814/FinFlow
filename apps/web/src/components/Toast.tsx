'use client';
import { useEffect } from 'react';

interface Props { message: string; kind?: 'success' | 'error'; onDismiss: () => void; durationMs?: number }

export function Toast({ message, kind = 'success', onDismiss, durationMs = 4000 }: Props) {
  useEffect(() => {
    const t = setTimeout(onDismiss, durationMs);
    return () => clearTimeout(t);
  }, [onDismiss, durationMs]);

  return (
    <div
      role={kind === 'error' ? 'alert' : 'status'}
      className={`fixed bottom-6 right-6 rounded-md px-4 py-3 text-sm text-white shadow-lg ${kind === 'error' ? 'bg-red-600' : 'bg-emerald-600'}`}
    >
      {message}
    </div>
  );
}
