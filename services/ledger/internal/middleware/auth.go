package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const userKey ctxKey = "userID"

func UserID(ctx context.Context) string {
	v, _ := ctx.Value(userKey).(string)
	return v
}

// Auth validates HS256 bearer tokens and stores the subject in the context.
func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if raw == "" {
				http.Error(w, `{"code":"unauthorized","message":"missing token"}`, http.StatusUnauthorized)
				return
			}
			claims := jwt.RegisteredClaims{}
			_, err := jwt.ParseWithClaims(raw, &claims, func(*jwt.Token) (any, error) {
				return []byte(secret), nil
			}, jwt.WithValidMethods([]string{"HS256"}))
			if err != nil || claims.Subject == "" {
				http.Error(w, `{"code":"unauthorized","message":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey, claims.Subject)))
		})
	}
}
