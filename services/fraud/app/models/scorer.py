from __future__ import annotations

import math
from collections import defaultdict
from datetime import datetime, timezone

from app.config import Settings
from app.models.features import AccountHistory, extract_features
from app.models.rules import RuleHit, evaluate_rules
from app.schemas import ScoreRequest, ScoreResponse, Signal


class FraudScorer:
    def __init__(self, settings: Settings):
        self.settings = settings
        self.histories: dict[str, AccountHistory] = defaultdict(AccountHistory)

    @staticmethod
    def _combine(hits: list[RuleHit]) -> float:
        """Noisy-OR combination keeps the score in [0, 1] as signals accumulate."""
        return 1.0 - math.prod(1.0 - h.weight for h in hits) if hits else 0.0

    def score(self, req: ScoreRequest, now: datetime | None = None) -> ScoreResponse:
        now = now or datetime.now(timezone.utc)
        history = self.histories[req.source_account]
        feats = extract_features(history, req.amount_minor, req.dest_account, now,
                                 self.settings.velocity_window_seconds)
        hits = evaluate_rules(feats, req.amount_minor, self.settings)
        score = round(self._combine(hits), 4)

        if score >= self.settings.block_threshold:
            action = "BLOCK"
        elif score >= self.settings.review_threshold:
            action = "REVIEW"
        else:
            action = "ALLOW"

        if action == "ALLOW":
            history.record(req.amount_minor, req.dest_account, now)
        return ScoreResponse(
            score=score,
            action=action,
            signals=[Signal(code=h.code, weight=h.weight, detail=h.detail) for h in hits],
            evaluated_at=now,
        )
