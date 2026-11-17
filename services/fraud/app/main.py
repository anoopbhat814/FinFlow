from fastapi import FastAPI

from app import __version__
from app.config import settings
from app.models import FraudScorer
from app.schemas import ScoreRequest, ScoreResponse

app = FastAPI(title="FinFlow Fraud Service", version=__version__)
scorer = FraudScorer(settings)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok", "version": __version__}


@app.post("/v1/score", response_model=ScoreResponse)
def score(req: ScoreRequest) -> ScoreResponse:
    return scorer.score(req)
