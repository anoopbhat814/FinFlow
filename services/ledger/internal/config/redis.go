package config

import (
	"net/url"
	"strconv"
)

// RedisConfig describes the connection used for idempotency keys and rate limiting.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// ParseRedisURL accepts URLs of the form redis://:password@host:6379/2.
func ParseRedisURL(raw string) (RedisConfig, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return RedisConfig{}, err
	}
	rc := RedisConfig{Addr: u.Host}
	if pw, ok := u.User.Password(); ok {
		rc.Password = pw
	}
	if p := u.Path; len(p) > 1 {
		if db, err := strconv.Atoi(p[1:]); err == nil {
			rc.DB = db
		}
	}
	return rc, nil
}
