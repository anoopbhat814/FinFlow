'use client';
import type { Card, Page } from '@finflow/types';
import { CardTile } from '@/components/CardTile';
import { api } from '@/lib/api-client';
import { useAsync } from '@/hooks/useAsync';

export default function CardsPage() {
  const { data, reload } = useAsync(() => api<Page<Card>>('/cards'));

  async function toggleFreeze(card: Card) {
    const action = card.status === 'FROZEN' ? 'unfreeze' : 'freeze';
    await api(`/cards/${card.id}/${action}`, { method: 'POST' });
    reload();
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      {data?.items.map((c) => <CardTile key={c.id} card={c} onToggleFreeze={toggleFreeze} />)}
    </div>
  );
}
