'use client';
import { useState, type FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api-client';
import { signIn } from '@/lib/auth';
import { Button } from '@/components/Button';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    try {
      const { token } = await api<{ token: string }>('/auth/login', { method: 'POST', body: { email, password } });
      if (!signIn(token)) throw new Error('Received an invalid session');
      router.push('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Sign-in failed');
    } finally {
      setBusy(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="mx-auto mt-24 max-w-sm space-y-4 rounded-xl bg-white p-6 shadow">
      <h1 className="text-xl font-semibold">Sign in to FinFlow</h1>
      <input type="email" required value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email" className="w-full rounded border p-2" />
      <input type="password" required value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" className="w-full rounded border p-2" />
      {error && <p role="alert" className="text-sm text-red-600">{error}</p>}
      <Button type="submit" loading={busy} className="w-full">Sign in</Button>
    </form>
  );
}
