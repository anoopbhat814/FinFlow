package handlers

import (
	"net/http"
	"strconv"

	"github.com/finflow/ledger/internal/middleware"
	"github.com/finflow/ledger/internal/repository"
	"github.com/finflow/ledger/internal/service"
)

type AccountHandler struct {
	repo   *repository.AccountRepo
	ledger *service.LedgerService
}

func NewAccountHandler(r *repository.AccountRepo, l *service.LedgerService) *AccountHandler {
	return &AccountHandler{repo: r, ledger: l}
}

func (h *AccountHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/accounts", h.list)
	mux.HandleFunc("GET /v1/accounts/{id}/entries", h.entries)
}

func (h *AccountHandler) list(w http.ResponseWriter, r *http.Request) {
	accts, err := h.repo.ListByOwner(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": accts, "nextCursor": nil})
}

func (h *AccountHandler) entries(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	entries, err := h.ledger.Statement(r.Context(), middleware.UserID(r.Context()), r.PathValue("id"), limit)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": entries, "nextCursor": nil})
}
