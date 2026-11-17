import { movingAverage, type DailyTotal } from '@/lib/analytics';

interface Props { data: DailyTotal[]; width?: number; height?: number }

export function SpendingChart({ data, width = 560, height = 160 }: Props) {
  if (data.length < 2) return <p className="text-sm text-slate-500">Not enough data yet.</p>;
  const values = data.map((d) => d.outflowMinor);
  const smooth = movingAverage(values, 3);
  const max = Math.max(...values, 1);
  const x = (i: number) => (i / (data.length - 1)) * width;
  const y = (v: number) => height - (v / max) * (height - 8) - 4;
  const path = smooth.map((v, i) => `${i === 0 ? 'M' : 'L'}${x(i).toFixed(1)},${y(v).toFixed(1)}`).join(' ');

  return (
    <svg viewBox={`0 0 ${width} ${height}`} role="img" aria-label="Daily outflow, 3-day average" className="w-full">
      {values.map((v, i) => (
        <rect key={data[i]!.date} x={x(i) - 3} y={y(v)} width={6} height={height - y(v)} className="fill-indigo-200" />
      ))}
      <path d={path} fill="none" strokeWidth={2} className="stroke-indigo-600" />
    </svg>
  );
}
