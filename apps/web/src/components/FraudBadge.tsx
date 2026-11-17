import type { FraudDecision, RiskLevel } from '@finflow/types';

const COLORS: Record<RiskLevel, string> = {
  LOW: 'bg-emerald-100 text-emerald-800',
  MEDIUM: 'bg-amber-100 text-amber-800',
  HIGH: 'bg-red-100 text-red-800',
};

export function FraudBadge({ decision }: { decision?: FraudDecision }) {
  if (!decision) return null;
  const title = decision.signals.map((s) => s.detail).join('\n') || 'No risk signals';
  return (
    <span title={title} className={`rounded-full px-2 py-0.5 text-xs font-medium ${COLORS[decision.level]}`}>
      {decision.level.toLowerCase()} risk · {(decision.score * 100).toFixed(0)}
    </span>
  );
}
