package domain

import (
	"errors"
	"time"
)

type RefundStatus string

const (
	RefundPending   RefundStatus = "PENDING"
	RefundSucceeded RefundStatus = "SUCCEEDED"
	RefundFailed    RefundStatus = "FAILED"
)

var (
	ErrTransferNotFound     = errors.New("transfer not found")
	ErrTransferNotRefundable = errors.New("transfer is not refundable")
	ErrRefundExceedsTransfer = errors.New("refund exceeds remaining refundable amount")
	ErrInvalidAmount        = errors.New("amount must be positive")
)

type Refund struct {
	ID          string       `json:"id"`
	TransferID  string       `json:"transferId"`
	AmountMinor int64        `json:"amountMinor"`
	Currency    string       `json:"currency"`
	Reason      string       `json:"reason,omitempty"`
	Status      RefundStatus `json:"status"`
	CreatedAt   time.Time    `json:"createdAt"`
	SettledAt   *time.Time   `json:"settledAt,omitempty"`
}

// RefundableRemaining returns how much of t may still be refunded given the
// refunds already recorded. Failed refunds do not consume the budget.
func RefundableRemaining(t Transfer, existing []Refund) int64 {
	left := t.AmountMinor
	for _, r := range existing {
		if r.Status != RefundFailed {
			left -= r.AmountMinor
		}
	}
	if left < 0 {
		return 0
	}
	return left
}

// ValidateRefund checks that amount can be refunded against t.
func ValidateRefund(t Transfer, existing []Refund, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if t.Status != TransferCompleted {
		return ErrTransferNotRefundable
	}
	if amount > RefundableRemaining(t, existing) {
		return ErrRefundExceedsTransfer
	}
	return nil
}

type DisputeStatus string

const (
	DisputeOpen             DisputeStatus = "OPEN"
	DisputeEvidenceRequired DisputeStatus = "EVIDENCE_REQUIRED"
	DisputeUnderReview      DisputeStatus = "UNDER_REVIEW"
	DisputeWon              DisputeStatus = "WON"
	DisputeLost             DisputeStatus = "LOST"
	DisputeWithdrawn        DisputeStatus = "WITHDRAWN"
)

var ErrInvalidDisputeTransition = errors.New("invalid dispute status transition")

var disputeTransitions = map[DisputeStatus][]DisputeStatus{
	DisputeOpen:             {DisputeEvidenceRequired, DisputeUnderReview, DisputeWithdrawn},
	DisputeEvidenceRequired: {DisputeUnderReview, DisputeWithdrawn},
	DisputeUnderReview:      {DisputeWon, DisputeLost, DisputeEvidenceRequired},
}

// CanTransition reports whether a dispute may move from one status to another.
func (from DisputeStatus) CanTransition(to DisputeStatus) bool {
	for _, next := range disputeTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

func (s DisputeStatus) Terminal() bool {
	return s == DisputeWon || s == DisputeLost || s == DisputeWithdrawn
}

type Dispute struct {
	ID         string        `json:"id"`
	TransferID string        `json:"transferId"`
	OpenedByID string        `json:"openedById"`
	Reason     string        `json:"reason"`
	Status     DisputeStatus `json:"status"`
	Details    string        `json:"details,omitempty"`
	OpenedAt   time.Time     `json:"openedAt"`
	ResolvedAt *time.Time    `json:"resolvedAt,omitempty"`
}

// Advance moves the dispute to the next status, stamping ResolvedAt on terminal states.
func (d *Dispute) Advance(to DisputeStatus, now time.Time) error {
	if !d.Status.CanTransition(to) {
		return ErrInvalidDisputeTransition
	}
	d.Status = to
	if to.Terminal() {
		d.ResolvedAt = &now
	}
	return nil
}
