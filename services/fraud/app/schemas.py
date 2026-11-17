from datetime import datetime
from typing import Literal

from pydantic import BaseModel, Field


class ScoreRequest(BaseModel):
    transfer_id: str
    source_account: str
    dest_account: str
    amount_minor: int = Field(gt=0)
    currency: str = Field(min_length=3, max_length=3)
    country: str | None = None


class Signal(BaseModel):
    code: str
    weight: float
    detail: str


class ScoreResponse(BaseModel):
    score: float = Field(ge=0, le=1)
    action: Literal["ALLOW", "REVIEW", "BLOCK"]
    signals: list[Signal]
    evaluated_at: datetime
