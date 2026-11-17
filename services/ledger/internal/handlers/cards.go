package handlers

import (
	"context"
	"net/http"

	"github.com/finflow/ledger/internal/middleware"
	"github.com/finflow/ledger/internal/repository"
	"github.com/finflow/ledger/internal/service"
)

type CardHandler struct {
	repo *repository.CardRepo
	svc  *service.CardService
}

func NewCardHandler(r *repository.CardRepo, s *service.CardService) *CardHandler {
	return &CardHandler{repo: r, svc: s}
}

func (h *CardHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/cards", h.list)
	mux.HandleFunc("POST /v1/cards/{id}/freeze", h.freeze)
	mux.HandleFunc("POST /v1/cards/{id}/unfreeze", h.unfreeze)
	mux.HandleFunc("POST /v1/cards/{id}/authorizations", h.authorize)
}

func (h *CardHandler) list(w http.ResponseWriter, r *http.Request) {
	cards, err := h.repo.ListByOwner(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": cards})
}

func (h *CardHandler) freeze(w http.ResponseWriter, r *http.Request) {
	h.toggle(w, r, h.svc.Freeze)
}

func (h *CardHandler) unfreeze(w http.ResponseWriter, r *http.Request) {
	h.toggle(w, r, h.svc.Unfreeze)
}

func (h *CardHandler) toggle(w http.ResponseWriter, r *http.Request, fn func(ctx context.Context, user, card string) error) {
	if err := fn(r.Context(), middleware.UserID(r.Context()), r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type authorizeRequest struct {
	Merchant    string `json:"merchant"`
	MCC         string `json:"mcc"`
	Currency    string `json:"currency"`
	AmountMinor int64  `json:"amountMinor"`
}

func (h *CardHandler) authorize(w http.ResponseWriter, r *http.Request) {
	var req authorizeRequest
	if err := decodeJSON(r, &req); err != nil || req.Merchant == "" || len(req.MCC) != 4 || len(req.Currency) != 3 {
		writeBadRequest(w, "merchant, 4-digit mcc and 3-letter currency are required")
		return
	}
	err := h.svc.Authorize(r.Context(), middleware.UserID(r.Context()), r.PathValue("id"), req.Merchant, req.MCC, req.Currency, req.AmountMinor)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"decision": "APPROVED"})
}
