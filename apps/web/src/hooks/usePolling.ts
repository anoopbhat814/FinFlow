import { useEffect, useRef } from 'react';

/** Polls `fn` every `intervalMs` until it returns true (done) or the component unmounts. */
export function usePolling(fn: () => Promise<boolean>, intervalMs: number, enabled = true) {
  const saved = useRef(fn);
  saved.current = fn;

  useEffect(() => {
    if (!enabled) return;
    let cancelled = false;
    let timer: ReturnType<typeof setTimeout>;
    const tick = async () => {
      const done = await saved.current().catch(() => false);
      if (!cancelled && !done) timer = setTimeout(tick, intervalMs);
    };
    timer = setTimeout(tick, intervalMs);
    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [intervalMs, enabled]);
}
