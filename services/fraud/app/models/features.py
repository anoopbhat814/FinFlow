from __future__ import annotations

import statistics
from dataclasses import dataclass, field
from datetime import datetime, timezone


@dataclass
class AccountHistory:
    """Rolling per-account history kept in memory (Redis-backed in production)."""

    amounts_minor: list[int] = field(default_factory=list)
    timestamps: list[datetime] = field(default_factory=list)
    known_destinations: set[str] = field(default_factory=set)

    def record(self, amount_minor: int, dest: str, at: datetime | None = None) -> None:
        self.amounts_minor.append(amount_minor)
        self.timestamps.append(at or datetime.now(timezone.utc))
        self.known_destinations.add(dest)


def extract_features(history: AccountHistory, amount_minor: int, dest: str, now: datetime,
                     window_seconds: int) -> dict[str, float]:
    recent = [t for t in history.timestamps if (now - t).total_seconds() <= window_seconds]
    mean = statistics.fmean(history.amounts_minor) if history.amounts_minor else 0.0
    stdev = statistics.pstdev(history.amounts_minor) if len(history.amounts_minor) > 1 else 0.0
    z = (amount_minor - mean) / stdev if stdev > 0 else 0.0
    return {
        "recent_count": float(len(recent)),
        "amount_zscore": z,
        "new_destination": 0.0 if dest in history.known_destinations else 1.0,
        "hour_of_day": float(now.hour),
        "history_size": float(len(history.amounts_minor)),
    }
