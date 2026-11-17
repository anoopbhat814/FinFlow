package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logging(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}
			w.Header().Set("X-Request-ID", reqID)
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			next.ServeHTTP(sw, r)
			log.Info("request", "id", reqID, "method", r.Method, "path", r.URL.Path,
				"status", sw.status, "duration_ms", time.Since(start).Milliseconds())
		})
	}
}
