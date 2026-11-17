from datetime import datetime, timedelta, timezone

from fastapi.testclient import TestClient

from app.config import Settings
from app.main import app
from app.models import FraudScorer
from app.schemas import ScoreRequest


def req(amount: int, dest: str = "acct_b") -> ScoreRequest:
    return ScoreRequest(transfer_id="t1", source_account="acct_a", dest_account=dest,
                        amount_minor=amount, currency="USD")


def test_small_transfer_allowed():
    res = FraudScorer(Settings()).score(req(1_500))
    assert res.action == "ALLOW" and res.score == 0.0


def test_velocity_burst_escalates_to_review():
    scorer = FraudScorer(Settings())
    now = datetime(2025, 3, 1, 12, tzinfo=timezone.utc)
    for i in range(6):
        res = scorer.score(req(1_000), now + timedelta(seconds=i))
    assert res.action in {"REVIEW", "BLOCK"}
    assert any(s.code == "VELOCITY" for s in res.signals)


def test_stacked_signals_block():
    scorer = FraudScorer(Settings(block_threshold=0.5))
    res = scorer.score(req(900_000), datetime(2025, 3, 1, 2, tzinfo=timezone.utc))
    assert res.action == "BLOCK"


def test_http_endpoint_contract():
    client = TestClient(app)
    r = client.post("/v1/score", json={"transfer_id": "x", "source_account": "a", "dest_account": "b",
                                       "amount_minor": 1200, "currency": "USD"})
    assert r.status_code == 200
    assert set(r.json()) == {"score", "action", "signals", "evaluated_at"}
