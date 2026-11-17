package domain

import (
	"testing"
	"time"
)

func TestValidateEndpointURL(t *testing.T) {
	tests := []struct {
		raw     string
		wantErr bool
	}{
		{"https://hooks.example.com/finflow", false},
		{"https://example.com:8443/x?y=1", false},
		{"http://example.com/hook", true},
		{"ftp://example.com", true},
		{"https://", true},
		{"not a url", true},
		{"", true},
	}
	for _, tc := range tests {
		t.Run(tc.raw, func(t *testing.T) {
			if err := ValidateEndpointURL(tc.raw); (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestSubscribes(t *testing.T) {
	tests := []struct {
		name  string
		ep    WebhookEndpoint
		event string
		want  bool
	}{
		{"exact match", WebhookEndpoint{Active: true, Events: []string{"transfer.completed"}}, "transfer.completed", true},
		{"wildcard", WebhookEndpoint{Active: true, Events: []string{"*"}}, "refund.succeeded", true},
		{"other event", WebhookEndpoint{Active: true, Events: []string{"transfer.completed"}}, "refund.succeeded", false},
		{"inactive", WebhookEndpoint{Active: false, Events: []string{"*"}}, "x", false},
		{"no events", WebhookEndpoint{Active: true}, "x", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.ep.Subscribes(tc.event); got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestVerifySignature(t *testing.T) {
	now := time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)
	body := []byte(`{"event":"transfer.completed"}`)
	good := Sign("s3cret", now, body)
	tests := []struct {
		name    string
		secret  string
		ts      time.Time
		body    []byte
		sig     string
		wantErr bool
	}{
		{"valid", "s3cret", now, body, good, false},
		{"valid within tolerance", "s3cret", now.Add(-4 * time.Minute), body, Sign("s3cret", now.Add(-4*time.Minute), body), false},
		{"wrong secret", "other", now, body, good, true},
		{"tampered body", "s3cret", now, []byte(`{}`), good, true},
		{"stale timestamp", "s3cret", now.Add(-10 * time.Minute), body, Sign("s3cret", now.Add(-10*time.Minute), body), true},
		{"future timestamp", "s3cret", now.Add(10 * time.Minute), body, Sign("s3cret", now.Add(10*time.Minute), body), true},
		{"empty signature", "s3cret", now, body, "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifySignature(tc.secret, tc.ts, now, 5*time.Minute, tc.body, tc.sig)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestBackoff(t *testing.T) {
	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{-1, 30 * time.Second}, {0, 30 * time.Second}, {1, 30 * time.Second},
		{2, time.Minute}, {3, 2 * time.Minute}, {6, 16 * time.Minute},
		{9, time.Hour}, {50, time.Hour},
	}
	for _, tc := range tests {
		if got := Backoff(tc.attempt); got != tc.want {
			t.Errorf("Backoff(%d) = %v, want %v", tc.attempt, got, tc.want)
		}
	}
}
