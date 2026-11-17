package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	tests := []struct {
		name       string
		inbound    string
		wantEchoed bool
	}{
		{"valid id is kept", "req-12345678", true},
		{"missing id is generated", "", false},
		{"too short is replaced", "abc", false},
		{"bad characters are replaced", "bad id with spaces!", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var seen string
			h := RequestID(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				seen = RequestIDFrom(r.Context())
			}))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.inbound != "" {
				req.Header.Set("X-Request-ID", tc.inbound)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if seen == "" || rec.Header().Get("X-Request-ID") != seen {
				t.Fatalf("context id %q, header %q", seen, rec.Header().Get("X-Request-ID"))
			}
			if tc.wantEchoed != (seen == tc.inbound) {
				t.Fatalf("id %q vs inbound %q", seen, tc.inbound)
			}
		})
	}
}
