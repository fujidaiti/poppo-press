package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	n, err := rr.ResponseWriter.Write(b)
	rr.size += n
	return n, err
}

func jsonLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rr, r)
		entry := map[string]any{
			"ts":       time.Now().UTC().Format(time.RFC3339Nano),
			"remoteIP": clientIP(r),
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   rr.status,
			"size":     rr.size,
			"duration": time.Since(start).Milliseconds(),
			"reqId":    middleware.GetReqID(r.Context()),
		}
		b, _ := json.Marshal(entry)
		_, _ = w.Write([]byte{})
		_ = b // avoid unused in case
		// Write to stderr via default logger by printing; but avoid interfering with response
		// Use standard library to print to stderr
		// We avoid importing log to keep output minimal; write via middleware.DefaultLogFormatter is not desired
		// Using println is acceptable for simple JSON logs
		println(string(b))
	})
}

func limitBody(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take first IP
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}
	if rip := r.Header.Get("X-Real-IP"); rip != "" {
		return rip
	}
	return r.RemoteAddr
}
