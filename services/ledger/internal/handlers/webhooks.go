package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/finflow/ledger/internal/domain"
	"github.com/finflow/ledger/internal/middleware"
	"github.com/finflow/ledger/internal/repository"
)

type WebhookHandler struct{ repo *repository.WebhookRepo }

func NewWebhookHandler(r *repository.WebhookRepo) *WebhookHandler { return &WebhookHandler{repo: r} }

func (h *WebhookHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/webhooks", h.list)
	mux.HandleFunc("POST /v1/webhooks", h.create)
}

type createWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

func (h *WebhookHandler) list(w http.ResponseWriter, r *http.Request) {
	eps, err := h.repo.ListEndpoints(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": eps})
}

// create returns the signing secret exactly once; it is never serialised again.
func (h *WebhookHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createWebhookRequest
	if err := decodeJSON(r, &req); err != nil || len(req.Events) == 0 || len(req.Events) > 20 {
		writeBadRequest(w, "url and between 1 and 20 events are required")
		return
	}
	if err := domain.ValidateEndpointURL(req.URL); err != nil {
		writeError(w, err)
		return
	}
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		writeError(w, err)
		return
	}
	secret := "whsec_" + hex.EncodeToString(buf)
	ep, err := h.repo.CreateEndpoint(r.Context(), domain.WebhookEndpoint{
		UserID: middleware.UserID(r.Context()), URL: req.URL, Secret: secret, Events: req.Events,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"endpoint": ep, "secret": secret})
}
