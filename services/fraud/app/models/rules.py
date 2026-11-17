from __future__ import annotations

from dataclasses import dataclass
from typing import Callable

from app.config import Settings


@dataclass(frozen=True)
class RuleHit:
    code: str
    weight: float
    detail: str


Rule = Callable[[dict[str, float], int, Settings], RuleHit | None]


def large_amount(f: dict[str, float], amount: int, s: Settings) -> RuleHit | None:
    if amount >= s.large_amount_minor:
        return RuleHit("LARGE_AMOUNT", 0.35, f"Amount {amount} exceeds {s.large_amount_minor}")
    return None


def velocity(f: dict[str, float], amount: int, s: Settings) -> RuleHit | None:
    if f["recent_count"] >= s.velocity_max_transfers:
        return RuleHit("VELOCITY", 0.4, f"{int(f['recent_count'])} transfers in the last window")
    return None


def new_destination_large(f: dict[str, float], amount: int, s: Settings) -> RuleHit | None:
    if f["new_destination"] and amount >= s.large_amount_minor // 5:
        return RuleHit("NEW_PAYEE", 0.2, "First transfer to this destination is sizeable")
    return None


def unusual_hour(f: dict[str, float], amount: int, s: Settings) -> RuleHit | None:
    if 1 <= f["hour_of_day"] <= 4 and amount > 50_000:
        return RuleHit("ODD_HOUR", 0.1, "Large transfer between 01:00 and 04:59 UTC")
    return None


def amount_outlier(f: dict[str, float], amount: int, s: Settings) -> RuleHit | None:
    if f["history_size"] >= 5 and f["amount_zscore"] > 3:
        return RuleHit("OUTLIER", 0.3, f"Amount is {f['amount_zscore']:.1f} sigma above normal")
    return None


RULES: list[Rule] = [large_amount, velocity, new_destination_large, unusual_hour, amount_outlier]


def evaluate_rules(features: dict[str, float], amount: int, settings: Settings) -> list[RuleHit]:
    return [hit for rule in RULES if (hit := rule(features, amount, settings))]
