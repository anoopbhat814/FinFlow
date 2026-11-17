package domain

import "time"

type TransferStatus string

const (
	TransferPending   TransferStatus = "PENDING"
	TransferInReview  TransferStatus = "IN_REVIEW"
	TransferCompleted TransferStatus = "COMPLETED"
	TransferFailed    TransferStatus = "FAILED"
	TransferReversed  TransferStatus = "REVERSED"
)

type Transfer struct {
	ID             string         `json:"id"`
	IdempotencyKey string         `json:"idempotencyKey"`
	SourceID       string         `json:"sourceAccountId"`
	DestinationID  string         `json:"destinationAccountId"`
	AmountMinor    int64          `json:"amountMinor"`
	Currency       string         `json:"currency"`
	Memo           string         `json:"memo,omitempty"`
	Status         TransferStatus `json:"status"`
	FraudScore     float64        `json:"fraudScore"`
	CreatedAt      time.Time      `json:"createdAt"`
	CompletedAt    *time.Time     `json:"completedAt,omitempty"`
}

type FraudVerdict struct {
	Score  float64 `json:"score"`
	Action string  `json:"action"` // ALLOW | REVIEW | BLOCK
}
