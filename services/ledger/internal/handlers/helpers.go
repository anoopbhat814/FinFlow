package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/finflow/ledger/internal/domain"
)

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{"not_found", err.Error()})
	case errors.Is(err, domain.ErrInsufficientFunds), errors.Is(err, domain.ErrAccountNotActive):
		writeJSON(w, http.StatusUnprocessableEntity, errorBody{"unprocessable", err.Error()})
	case errors.Is(err, domain.ErrTransferNotFound), errors.Is(err, domain.ErrCardNotFound), errors.Is(err, domain.ErrInvoiceNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{"not_found", err.Error()})
	case errors.Is(err, domain.ErrInvalidAmount), errors.Is(err, domain.ErrInvalidPeriod), errors.Is(err, domain.ErrInvalidWebhookURL):
		writeJSON(w, http.StatusBadRequest, errorBody{"invalid_request", err.Error()})
	case errors.Is(err, domain.ErrTransferNotRefundable), errors.Is(err, domain.ErrRefundExceedsTransfer),
		errors.Is(err, domain.ErrCardNotActive), errors.Is(err, domain.ErrCardExpired),
		errors.Is(err, domain.ErrLimitExceeded), errors.Is(err, domain.ErrInvoiceLocked):
		writeJSON(w, http.StatusUnprocessableEntity, errorBody{"unprocessable", err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorBody{"internal", "internal error"})
	}
}

// decodeJSON reads a bounded JSON body and rejects unknown fields.
func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func writeBadRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, errorBody{"invalid_request", msg})
}
