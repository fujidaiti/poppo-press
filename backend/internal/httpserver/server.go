// Package httpserver provides the HTTP mux and common middleware for the API
// server, including simple health and version endpoints.
package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"database/sql"

	"github.com/fujidaiti/poppo-press/backend/internal/auth"
	"github.com/fujidaiti/poppo-press/backend/internal/db"
	"github.com/fujidaiti/poppo-press/backend/internal/version"
)

// Server wraps a chi router that exposes core service endpoints and middleware.
// It is responsible for wiring health and version routes and returning the
// http.Handler used by the HTTP server.
type Server struct {
	mux *chi.Mux
	db  *sql.DB
}

// New constructs a Server with standard middleware (RealIP, RequestID, Logger,
// Recoverer) and registers the /health and /version endpoints.
func New(database *sql.DB) *Server {
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

	// v1 routes
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
				type req struct{ Username, Password, DeviceName string }
				var body req
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					writeError(w, http.StatusBadRequest, "bad_request", "invalid json")
					return
				}
				if body.Username == "" || body.Password == "" || body.DeviceName == "" {
					writeError(w, http.StatusBadRequest, "validation_failed", "missing fields")
					return
				}
				phc, err := db.GetUserPasswordHash(r.Context(), database, body.Username)
				if err != nil || phc == "" {
					writeError(w, http.StatusUnauthorized, "unauthorized", "invalid credentials")
					return
				}
				ok, _ := auth.VerifyPassword(body.Password, phc)
				if !ok {
					writeError(w, http.StatusUnauthorized, "unauthorized", "invalid credentials")
					return
				}
				token, err := auth.GenerateToken()
				if err != nil {
					writeError(w, http.StatusInternalServerError, "internal", "failed to generate token")
					return
				}
				id, err := db.CreateOrUpdateDeviceToken(r.Context(), database, body.DeviceName, auth.HashToken(token))
				if err != nil {
					writeError(w, http.StatusInternalServerError, "internal", "failed to persist token")
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{"token": token, "deviceId": id})
			})
			r.With(authMiddleware(database)).Post("/logout", func(w http.ResponseWriter, r *http.Request) {
				devID := r.Context().Value(ctxDeviceID{}).(int64)
				if err := db.RevokeDeviceToken(r.Context(), database, devID); err != nil {
					writeError(w, http.StatusInternalServerError, "internal", "failed to revoke")
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})
		})

		// protected test route for M2 DoD
		r.With(authMiddleware(database)).Get("/protected/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"pong": "ok"})
		})

		// M3 Sources API
		registerSourcesRoutes(database, r)

		// M5 Editions API
		registerEditionRoutes(database, r)

		// M6 Articles API
		registerArticleRoutes(database, r)

		// M7 Read Later API
		registerReadLaterRoutes(database, r)
	})

	return &Server{mux: r, db: database}
}

// Handler returns the underlying http.Handler to serve requests.
func (s *Server) Handler() http.Handler { return s.mux }

// error helpers and auth middleware
type ctxDeviceID struct{}

func writeError(w http.ResponseWriter, code int, errCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"code": errCode, "message": msg}})
}

func authMiddleware(database *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if authz == "" {
				writeError(w, http.StatusUnauthorized, "unauthorized", "missing token")
				return
			}
			const prefix = "Bearer "
			if len(authz) <= len(prefix) || authz[:len(prefix)] != prefix {
				writeError(w, http.StatusUnauthorized, "unauthorized", "invalid token format")
				return
			}
			token := authz[len(prefix):]
			devID, err := db.LookupDeviceIdByToken(r.Context(), database, token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "invalid token")
				return
			}
			ctx := contextWithDeviceID(r.Context(), devID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func contextWithDeviceID(ctx context.Context, id int64) context.Context {
	type key = ctxDeviceID
	return context.WithValue(ctx, key{}, id)
}
