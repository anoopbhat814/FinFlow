import { ApiError, type ApiErrorBody } from '@finflow/types';
import { useSession } from '@/store/session';

const BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';

export interface RequestOptions extends Omit<RequestInit, 'body'> {
  body?: unknown;
  idempotencyKey?: string;
}

export async function api<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { token } = useSession.getState();
  const headers = new Headers(opts.headers);
  headers.set('Accept', 'application/json');
  if (opts.body !== undefined) headers.set('Content-Type', 'application/json');
  if (token) headers.set('Authorization', `Bearer ${token}`);
  if (opts.idempotencyKey) headers.set('Idempotency-Key', opts.idempotencyKey);

  const res = await fetch(`${BASE}/v1${path}`, {
    ...opts,
    headers,
    body: opts.body === undefined ? undefined : JSON.stringify(opts.body),
  });

  if (res.status === 401) useSession.getState().clear();
  if (!res.ok) {
    const body = (await res.json().catch(() => ({ code: 'unknown', message: res.statusText }))) as ApiErrorBody;
    throw new ApiError(res.status, body);
  }
  return res.status === 204 ? (undefined as T) : ((await res.json()) as T);
}

export const newIdempotencyKey = (): string => crypto.randomUUID();
