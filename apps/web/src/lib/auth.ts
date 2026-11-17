import { useSession } from '@/store/session';

interface JwtPayload { sub: string; exp: number }

export function decodeToken(token: string): JwtPayload | null {
  try {
    const [, payload] = token.split('.');
    if (!payload) return null;
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/'))) as JwtPayload;
  } catch {
    return null;
  }
}

export function isExpired(token: string, nowMs = Date.now(), skewSec = 30): boolean {
  const p = decodeToken(token);
  return !p || p.exp * 1000 <= nowMs + skewSec * 1000;
}

export function signIn(token: string): boolean {
  const p = decodeToken(token);
  if (!p || isExpired(token)) return false;
  useSession.getState().setSession(token, p.sub);
  return true;
}
