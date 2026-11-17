package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/finflow/ledger/internal/domain"
)

type WebhookStore interface {
	ListEndpoints(ctx context.Context, userID string) ([]domain.WebhookEndpoint, error)
	EnqueueDelivery(ctx context.Context, endpointID, event string, payload []byte) error
	DueDeliveries(ctx context.Context, now time.Time, limit int) ([]domain.WebhookDelivery, error)
	RecordAttempt(ctx context.Context, id, status string, code int, next *time.Time) error
}

type WebhookService struct {
	store       WebhookStore
	client      *http.Client
	maxAttempts int
	endpoints   func(ctx context.Context, id string) (*domain.WebhookEndpoint, error)
	now         func() time.Time
}

// NewWebhookService wires a dispatcher. lookup resolves an endpoint (with its
// secret) by id; it is injected so tests need no database.
func NewWebhookService(store WebhookStore, lookup func(context.Context, string) (*domain.WebhookEndpoint, error)) *WebhookService {
	return &WebhookService{
		store:       store,
		client:      &http.Client{Timeout: 5 * time.Second},
		maxAttempts: 8,
		endpoints:   lookup,
		now:         time.Now,
	}
}

// Publish fans an event out to every active endpoint the user subscribed.
func (s *WebhookService) Publish(ctx context.Context, userID, event string, data any) (int, error) {
	eps, err := s.store.ListEndpoints(ctx, userID)
	if err != nil {
		return 0, err
	}
	body, err := json.Marshal(map[string]any{"event": event, "data": data, "createdAt": s.now().UTC()})
	if err != nil {
		return 0, err
	}
	n := 0
	for _, ep := range eps {
		if !ep.Subscribes(event) {
			continue
		}
		if err := s.store.EnqueueDelivery(ctx, ep.ID, event, body); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// DeliverDue sends pending deliveries, scheduling retries with exponential
// backoff and marking them FAILED after maxAttempts.
func (s *WebhookService) DeliverDue(ctx context.Context, batch int) (delivered int, err error) {
	due, err := s.store.DueDeliveries(ctx, s.now(), batch)
	if err != nil {
		return 0, err
	}
	for _, d := range due {
		ep, err := s.endpoints(ctx, d.EndpointID)
		if err != nil {
			return delivered, err
		}
		code := s.send(ctx, ep, d)
		attempt := d.Attempts + 1
		switch {
		case code >= 200 && code < 300:
			delivered++
			err = s.store.RecordAttempt(ctx, d.ID, "DELIVERED", code, nil)
		case attempt >= s.maxAttempts:
			err = s.store.RecordAttempt(ctx, d.ID, "FAILED", code, nil)
		default:
			next := s.now().Add(domain.Backoff(attempt))
			err = s.store.RecordAttempt(ctx, d.ID, "PENDING", code, &next)
		}
		if err != nil {
			return delivered, err
		}
	}
	return delivered, nil
}

func (s *WebhookService) send(ctx context.Context, ep *domain.WebhookEndpoint, d domain.WebhookDelivery) int {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep.URL, bytes.NewReader(d.Payload))
	if err != nil {
		return 0
	}
	ts := s.now()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-FinFlow-Event", d.Event)
	req.Header.Set("X-FinFlow-Timestamp", strconv.FormatInt(ts.Unix(), 10))
	req.Header.Set("X-FinFlow-Signature", domain.Sign(ep.Secret, ts, d.Payload))
	resp, err := s.client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
