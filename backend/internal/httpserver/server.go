// Package httpserver provides the HTTP mux and common middleware for the API
// server, including simple health and version endpoints.
package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/fujidaiti/poppo-press/backend/internal/version"
)

// Server wraps a chi router that exposes core service endpoints and middleware.
// It is responsible for wiring health and version routes and returning the
// http.Handler used by the HTTP server.
type Server struct {
	mux *chi.Mux
}

// New constructs a Server with standard middleware (RealIP, RequestID, Logger,
// Recoverer) and registers the /health and /version endpoints.
func New() *Server {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"version": version.Version,
			"commit":  version.Commit,
			"date":    version.Date,
		})
	})

	return &Server{mux: r}
}

// Handler returns the underlying http.Handler to serve requests.
func (s *Server) Handler() http.Handler { return s.mux }
