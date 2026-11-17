package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/finflow/ledger/internal/domain"
)

// HTTPFraudClient calls the Python fraud service.
type HTTPFraudClient struct {
	baseURL string
	http    *http.Client
}

func NewHTTPFraudClient(baseURL string, timeout time.Duration) *HTTPFraudClient {
	return &HTTPFraudClient{baseURL: baseURL, http: &http.Client{Timeout: timeout}}
}

func (c *HTTPFraudClient) Check(ctx context.Context, t domain.Transfer) (domain.FraudVerdict, error) {
	payload, _ := json.Marshal(map[string]any{
		"transfer_id":    t.ID,
		"source_account": t.SourceID,
		"dest_account":   t.DestinationID,
		"amount_minor":   t.AmountMinor,
		"currency":       t.Currency,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/score", bytes.NewReader(payload))
	if err != nil {
		return domain.FraudVerdict{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.http.Do(req)
	if err != nil {
		return domain.FraudVerdict{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return domain.FraudVerdict{}, fmt.Errorf("fraud service returned %d", res.StatusCode)
	}
	var v domain.FraudVerdict
	return v, json.NewDecoder(res.Body).Decode(&v)
}
