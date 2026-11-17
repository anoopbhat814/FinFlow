from app.config import Settings
from app.models.rules import evaluate_rules

S = Settings()


def feats(**over):
    base = {"recent_count": 0.0, "amount_zscore": 0.0, "new_destination": 0.0,
            "hour_of_day": 12.0, "history_size": 0.0}
    base.update(over)
    return base


def test_no_hits_for_small_routine_transfer():
    assert evaluate_rules(feats(), 2_000, S) == []


def test_large_amount_triggers():
    codes = {h.code for h in evaluate_rules(feats(), 600_000, S)}
    assert "LARGE_AMOUNT" in codes


def test_velocity_triggers_at_threshold():
    codes = {h.code for h in evaluate_rules(feats(recent_count=5.0), 1_000, S)}
    assert codes == {"VELOCITY"}


def test_outlier_requires_history():
    assert evaluate_rules(feats(amount_zscore=4.0, history_size=2.0), 1_000, S) == []
    hits = evaluate_rules(feats(amount_zscore=4.0, history_size=10.0), 1_000, S)
    assert [h.code for h in hits] == ["OUTLIER"]
