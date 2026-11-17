export default function Loading() {
  return (
    <div className="animate-pulse space-y-4" aria-busy="true" aria-label="Loading your accounts">
      <div className="h-7 w-40 rounded bg-slate-200" />
      <div className="grid gap-4 md:grid-cols-2">
        {Array.from({ length: 2 }, (_, i) => (
          <div key={i} className="h-28 rounded-xl bg-slate-200" />
        ))}
      </div>
    </div>
  );
}
