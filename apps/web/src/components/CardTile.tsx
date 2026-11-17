import { remainingLimit, type Card as CardModel } from '@finflow/types';
import { formatMinor } from '@/lib/format';
import { Button } from './Button';

interface Props {
  card: CardModel;
  onToggleFreeze: (card: CardModel) => void;
}

export function CardTile({ card, onToggleFreeze }: Props) {
  const used = card.monthlyLimitMinor ? (card.spentThisMonthMinor / card.monthlyLimitMinor) * 100 : 0;
  return (
    <div className="rounded-xl bg-gradient-to-br from-indigo-600 to-slate-800 p-5 text-white">
      <p className="text-xs uppercase opacity-70">{card.type} · {card.brand}</p>
      <p className="mt-6 font-mono text-lg tracking-widest">•••• •••• •••• {card.last4}</p>
      <p className="mt-2 text-xs opacity-80">
        Exp {String(card.expMonth).padStart(2, '0')}/{String(card.expYear).slice(-2)} · {formatMinor(remainingLimit(card), card.currency)} left
      </p>
      <div className="mt-3 h-1.5 rounded bg-white/20">
        <div className="h-1.5 rounded bg-white" style={{ width: `${Math.min(100, used)}%` }} />
      </div>
      <Button variant="secondary" className="mt-4" onClick={() => onToggleFreeze(card)}>
        {card.status === 'FROZEN' ? 'Unfreeze' : 'Freeze'}
      </Button>
    </div>
  );
}
