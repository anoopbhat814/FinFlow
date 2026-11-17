// Package idempotency stores request results so retries are safe.
package idempotency

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrInFlight = errors.New("idempotency: request already in flight")

type record struct {
	response []byte
	done     bool
	expires  time.Time
}

// Store is an in-memory implementation; production deployments use Redis.
type Store struct {
	mu   sync.Mutex
	ttl  time.Duration
	data map[string]*record
}

func NewStore(ttl time.Duration) *Store {
	return &Store{ttl: ttl, data: make(map[string]*record)}
}

// Begin reserves key. If a completed response exists it is returned with ok=true.
func (s *Store) Begin(_ context.Context, key string) (resp []byte, ok bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if r, found := s.data[key]; found && r.expires.After(now) {
		if !r.done {
			return nil, false, ErrInFlight
		}
		return r.response, true, nil
	}
	s.data[key] = &record{expires: now.Add(s.ttl)}
	return nil, false, nil
}

func (s *Store) Complete(_ context.Context, key string, response []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.data[key]; ok {
		r.response, r.done = response, true
	}
}

func (s *Store) Abort(_ context.Context, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
