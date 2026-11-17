package handlers

import (
	"net/http"

	"github.com/finflow/ledger/internal/middleware"
	"github.com/finflow/ledger/internal/service"
)

type InvoiceHandler struct{ svc *service.InvoiceService }

func NewInvoiceHandler(s *service.InvoiceService) *InvoiceHandler { return &InvoiceHandler{svc: s} }

func (h *InvoiceHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/invoices", h.list)
	mux.HandleFunc("GET /v1/invoices/{id}", h.get)
	mux.HandleFunc("POST /v1/invoices/{id}/pay", h.pay)
}

func (h *InvoiceHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context(), middleware.UserID(r.Context()), r.URL.Query().Get("status"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h *InvoiceHandler) get(w http.ResponseWriter, r *http.Request) {
	inv, err := h.svc.Get(r.Context(), middleware.UserID(r.Context()), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, inv)
}

func (h *InvoiceHandler) pay(w http.ResponseWriter, r *http.Request) {
	inv, err := h.svc.MarkPaid(r.Context(), middleware.UserID(r.Context()), r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, inv)
}
