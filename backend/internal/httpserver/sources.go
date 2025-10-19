package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mmcdole/gofeed"

	"database/sql"

	"github.com/fujidaiti/poppo-press/backend/internal/db"
)

func registerSourcesRoutes(database *sql.DB, r chi.Router) {
	r.With(authMiddleware(database)).Route("/sources", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			type req struct{ URL string }
			var body req
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid json")
				return
			}
			if !isValidHTTPURL(body.URL) {
				writeError(w, http.StatusBadRequest, "validation_failed", "invalid url")
				return
			}
			title, etag, lastMod, err := probeFeed(r.Context(), body.URL)
			if err != nil {
				writeError(w, http.StatusBadRequest, "validation_failed", "unreachable or invalid feed")
				return
			}
			id, err := db.CreateSource(r.Context(), database, body.URL, title, etag, lastMod)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "failed to persist source")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
		})

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			rows, err := db.ListSources(r.Context(), database)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "failed to list sources")
				return
			}
			type out struct {
				ID        int64  `json:"id"`
				URL       string `json:"url"`
				Title     string `json:"title"`
				CreatedAt string `json:"createdAt"`
			}
			resp := make([]out, 0, len(rows))
			for _, r := range rows {
				resp = append(resp, out{ID: r.ID, URL: r.URL, Title: r.Title, CreatedAt: r.CreatedAt})
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})

		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, "id")
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid id")
				return
			}
			ok, err := db.DeleteSource(r.Context(), database, id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "failed to delete")
				return
			}
			if !ok {
				writeError(w, http.StatusNotFound, "not_found", "source not found")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})
}

func isValidHTTPURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

func probeFeed(ctx context.Context, rawURL string) (string, string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", "", "", err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", "", errors.New("bad status")
	}
	parser := gofeed.NewParser()
	feed, err := parser.Parse(resp.Body)
	if err != nil {
		return "", "", "", err
	}
	etag := resp.Header.Get("ETag")
	lastMod := resp.Header.Get("Last-Modified")
	return feed.Title, etag, lastMod, nil
}
