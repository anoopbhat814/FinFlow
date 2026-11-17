export type RiskLevel = 'LOW' | 'MEDIUM' | 'HIGH';
export type FraudAction = 'ALLOW' | 'REVIEW' | 'BLOCK';

export interface FraudSignal {
  code: string;
  weight: number;
  detail: string;
}

export interface FraudDecision {
  score: number;
  level: RiskLevel;
  action: FraudAction;
  signals: FraudSignal[];
  evaluatedAt: string;
}

export function levelFromScore(score: number): RiskLevel {
  if (score >= 0.75) return 'HIGH';
  if (score >= 0.4) return 'MEDIUM';
  return 'LOW';
}
