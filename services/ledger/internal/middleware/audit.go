package middleware

import (
	"context"
	"net"
	"net/http"

	"github.com/finflow/ledger/internal/domain"
)

// Recorder is satisfied by service.AuditService.
type Recorder interface {
	Record(ctx context.Context, e domain.AuditEvent)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Audit records every successful state-changing request (POST, PUT, PATCH,
// DELETE) with the caller, path, status and request id.
func Audit(rec Recorder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, r)
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			default:
				return
			}
			if sw.status >= 400 {
				return
			}
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			rec.Record(r.Context(), domain.AuditEvent{
				ActorID:    UserID(r.Context()),
				Action:     "http." + r.Method,
				EntityType: "route",
				EntityID:   r.URL.Path,
				IP:         ip,
				Metadata: map[string]string{
					"status":    http.StatusText(sw.status),
					"requestId": RequestIDFrom(r.Context()),
				},
			})
		})
	}
}
