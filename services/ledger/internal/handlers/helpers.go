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
	default:
		writeJSON(w, http.StatusInternalServerError, errorBody{"internal", "internal error"})
	}
}
