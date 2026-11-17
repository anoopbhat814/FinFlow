package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/finflow/ledger/internal/domain"
	"github.com/finflow/ledger/internal/service"
	"github.com/finflow/ledger/pkg/idempotency"
)

type TransferHandler struct {
	svc   *service.TransferService
	idem  *idempotency.Store
}

func NewTransferHandler(s *service.TransferService, i *idempotency.Store) *TransferHandler {
	return &TransferHandler{svc: s, idem: i}
}

func (h *TransferHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/transfers", h.create)
}

type createTransferReq struct {
	SourceID      string `json:"sourceAccountId"`
	DestinationID string `json:"destinationAccountId"`
	AmountMinor   int64  `json:"amountMinor"`
	Currency      string `json:"currency"`
	Memo          string `json:"memo"`
}

func (h *TransferHandler) create(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		writeJSON(w, http.StatusBadRequest, errorBody{"bad_request", "Idempotency-Key header is required"})
		return
	}
	if cached, ok, err := h.idem.Begin(r.Context(), key); err == idempotency.ErrInFlight {
		writeJSON(w, http.StatusConflict, errorBody{"in_flight", err.Error()})
		return
	} else if ok {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Idempotent-Replay", "true")
		_, _ = w.Write(cached)
		return
	}

	var req createTransferReq
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&req); err != nil || req.AmountMinor <= 0 {
		h.idem.Abort(r.Context(), key)
		writeJSON(w, http.StatusBadRequest, errorBody{"bad_request", "invalid transfer payload"})
		return
	}
	t, err := h.svc.Execute(r.Context(), domain.Transfer{
		IdempotencyKey: key, SourceID: req.SourceID, DestinationID: req.DestinationID,
		AmountMinor: req.AmountMinor, Currency: req.Currency, Memo: req.Memo,
	})
	if err != nil {
		h.idem.Abort(r.Context(), key)
		writeError(w, err)
		return
	}
	body, _ := json.Marshal(t)
	h.idem.Complete(r.Context(), key, body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(body)
}
