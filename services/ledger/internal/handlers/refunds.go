package handlers

import (
	"net/http"

	"github.com/finflow/ledger/internal/middleware"
	"github.com/finflow/ledger/internal/service"
)

type RefundHandler struct{ svc *service.RefundService }

func NewRefundHandler(s *service.RefundService) *RefundHandler { return &RefundHandler{svc: s} }

func (h *RefundHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/transfers/{id}/refunds", h.create)
}

type refundRequest struct {
	AmountMinor int64  `json:"amountMinor"`
	Reason      string `json:"reason"`
}

func (h *RefundHandler) create(w http.ResponseWriter, r *http.Request) {
	var req refundRequest
	if err := decodeJSON(r, &req); err != nil {
		writeBadRequest(w, "malformed json body")
		return
	}
	if req.AmountMinor < 0 || len(req.Reason) > 500 {
		writeBadRequest(w, "invalid amount or reason")
		return
	}
	ref, err := h.svc.Refund(r.Context(), middleware.UserID(r.Context()), r.PathValue("id"), req.AmountMinor, req.Reason)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, ref)
}
