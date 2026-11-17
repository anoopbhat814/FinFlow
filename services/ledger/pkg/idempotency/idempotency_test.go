package idempotency

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBeginReturnsStoredResponse(t *testing.T) {
	s := NewStore(time.Minute)
	ctx := context.Background()

	if _, ok, err := s.Begin(ctx, "k1"); ok || err != nil {
		t.Fatalf("first Begin: ok=%v err=%v", ok, err)
	}
	s.Complete(ctx, "k1", []byte(`{"id":"t_1"}`))

	resp, ok, err := s.Begin(ctx, "k1")
	if err != nil || !ok {
		t.Fatalf("second Begin: ok=%v err=%v", ok, err)
	}
	if string(resp) != `{"id":"t_1"}` {
		t.Fatalf("unexpected response %q", resp)
	}
}

func TestBeginInFlight(t *testing.T) {
	s := NewStore(time.Minute)
	ctx := context.Background()
	_, _, _ = s.Begin(ctx, "k2")
	if _, _, err := s.Begin(ctx, "k2"); !errors.Is(err, ErrInFlight) {
		t.Fatalf("err = %v, want ErrInFlight", err)
	}
}

func TestAbortReleasesKey(t *testing.T) {
	s := NewStore(time.Minute)
	ctx := context.Background()
	_, _, _ = s.Begin(ctx, "k3")
	s.Abort(ctx, "k3")
	if _, ok, err := s.Begin(ctx, "k3"); ok || err != nil {
		t.Fatalf("Begin after Abort: ok=%v err=%v", ok, err)
	}
}

func TestExpiredKeyIsReusable(t *testing.T) {
	s := NewStore(time.Millisecond)
	ctx := context.Background()
	_, _, _ = s.Begin(ctx, "k4")
	time.Sleep(5 * time.Millisecond)
	if _, _, err := s.Begin(ctx, "k4"); err != nil {
		t.Fatalf("expired key should be reusable, got %v", err)
	}
}
