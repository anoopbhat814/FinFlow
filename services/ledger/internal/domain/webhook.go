package domain

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrInvalidWebhookURL = errors.New("webhook url must be https")
	ErrBadSignature      = errors.New("invalid webhook signature")
)

type WebhookEndpoint struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	URL       string    `json:"url"`
	Secret    string    `json:"-"`
	Events    []string  `json:"events"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
}

type WebhookDelivery struct {
	ID            string     `json:"id"`
	EndpointID    string     `json:"endpointId"`
	Event         string     `json:"event"`
	Payload       []byte     `json:"payload"`
	Status        string     `json:"status"`
	Attempts      int        `json:"attempts"`
	ResponseCode  *int       `json:"responseCode,omitempty"`
	NextAttemptAt *time.Time `json:"nextAttemptAt,omitempty"`
}

// ValidateEndpointURL accepts only absolute https URLs with a host.
func ValidateEndpointURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return ErrInvalidWebhookURL
	}
	return nil
}

// Subscribes reports whether the endpoint is active and listens for event.
func (w WebhookEndpoint) Subscribes(event string) bool {
	if !w.Active {
		return false
	}
	for _, e := range w.Events {
		if e == event || e == "*" {
			return true
		}
	}
	return false
}

// Sign returns the hex HMAC-SHA256 over "<unix>.<body>" so receivers can
// reject replays by checking the timestamp.
func Sign(secret string, ts time.Time, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(ts.Unix(), 10)))
	mac.Write([]byte("."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature checks sig in constant time and enforces a tolerance window.
func VerifySignature(secret string, ts, now time.Time, tolerance time.Duration, body []byte, sig string) error {
	if d := now.Sub(ts); d > tolerance || d < -tolerance {
		return ErrBadSignature
	}
	if !hmac.Equal([]byte(Sign(secret, ts, body)), []byte(sig)) {
		return ErrBadSignature
	}
	return nil
}

// Backoff returns the delay before retry number attempt (1-based), doubling
// from 30s and capped at one hour.
func Backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := 30 * time.Second
	for i := 1; i < attempt && d < time.Hour; i++ {
		d *= 2
	}
	if d > time.Hour {
		d = time.Hour
	}
	return d
}
