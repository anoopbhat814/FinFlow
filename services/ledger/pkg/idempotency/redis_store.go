package idempotency

import (
	"context"
	"time"
)

// Cache is the subset of a Redis client the store needs. It keeps this package
// free of a direct client dependency and makes the store easy to fake.
type Cache interface {
	SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error)
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

const inFlightMarker = "__in_flight__"

// RedisStore shares idempotency state between ledger replicas.
type RedisStore struct {
	cache Cache
	ttl   time.Duration
}

func NewRedisStore(c Cache, ttl time.Duration) *RedisStore {
	return &RedisStore{cache: c, ttl: ttl}
}

func (s *RedisStore) Begin(ctx context.Context, key string) ([]byte, bool, error) {
	acquired, err := s.cache.SetNX(ctx, key, []byte(inFlightMarker), s.ttl)
	if err != nil {
		return nil, false, err
	}
	if acquired {
		return nil, false, nil
	}
	val, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, false, err
	}
	if string(val) == inFlightMarker {
		return nil, false, ErrInFlight
	}
	return val, true, nil
}

func (s *RedisStore) Complete(ctx context.Context, key string, response []byte) {
	_ = s.cache.Set(ctx, key, response, s.ttl)
}

func (s *RedisStore) Abort(ctx context.Context, key string) {
	_ = s.cache.Del(ctx, key)
}
