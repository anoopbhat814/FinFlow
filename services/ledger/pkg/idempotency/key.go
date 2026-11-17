package idempotency

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const maxKeyLength = 128

var ErrInvalidKey = errors.New("idempotency: key must be 1-128 printable characters")

// ScopedKey namespaces a client-supplied key by caller and route so two
// customers can never collide on the same Idempotency-Key header value.
func ScopedKey(userID, route, key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" || len(key) > maxKeyLength {
		return "", ErrInvalidKey
	}
	for _, r := range key {
		if r < 0x20 || r > 0x7e {
			return "", ErrInvalidKey
		}
	}
	sum := sha256.Sum256([]byte(userID + "|" + route + "|" + key))
	return hex.EncodeToString(sum[:]), nil
}
