package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/finflow/ledger/internal/middleware"
	"github.com/finflow/ledger/internal/service"
)

type StatementHandler struct{ svc *service.StatementService }

func NewStatementHandler(s *service.StatementService) *StatementHandler {
	return &StatementHandler{svc: s}
}

func (h *StatementHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/accounts/{id}/statements/{year}/{month}", h.monthly)
}

func (h *StatementHandler) monthly(w http.ResponseWriter, r *http.Request) {
	year, err1 := strconv.Atoi(r.PathValue("year"))
	month, err2 := strconv.Atoi(r.PathValue("month"))
	if err1 != nil || err2 != nil || year < 2000 || year > 2100 || month < 1 || month > 12 {
		writeBadRequest(w, "year and month must be valid")
		return
	}
	st, err := h.svc.Monthly(r.Context(), middleware.UserID(r.Context()), r.PathValue("id"), year, time.Month(month))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, st)
}
