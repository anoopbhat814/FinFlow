from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    review_threshold: float = 0.4
    block_threshold: float = 0.8
    large_amount_minor: int = 500_000
    velocity_window_seconds: int = 600
    velocity_max_transfers: int = 5

    model_config = {"env_prefix": "FRAUD_"}


settings = Settings()
