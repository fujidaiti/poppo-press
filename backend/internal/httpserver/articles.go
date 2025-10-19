package httpserver

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fujidaiti/poppo-press/backend/internal/db"
)

func registerArticleRoutes(database *sql.DB, r chi.Router) {
	r.With(authMiddleware(database)).Route("/articles", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			devID := r.Context().Value(ctxDeviceID{}).(int64)
			readState := r.URL.Query().Get("readState")
			if readState != "read" && readState != "unread" {
				readState = "all"
			}
			list, err := db.ListArticles(r.Context(), database, devID, readState)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "list fail")
				return
			}
			type out struct {
				ID           int64  `json:"id"`
				SourceID     int64  `json:"sourceId"`
				CanonicalURL string `json:"canonicalUrl"`
				Title        string `json:"title"`
				Summary      string `json:"summary"`
				Author       string `json:"author"`
				PublishedAt  string `json:"publishedAt"`
				IsRead       bool   `json:"isRead"`
			}
			resp := make([]out, 0, len(list))
			for _, a := range list {
				resp = append(resp, out{ID: a.ID, SourceID: a.SourceID, CanonicalURL: a.CanonicalURL, Title: a.Title, Summary: a.Summary, Author: a.Author, PublishedAt: a.PublishedAt, IsRead: a.IsRead})
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			devID := r.Context().Value(ctxDeviceID{}).(int64)
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid id")
				return
			}
			a, err := db.GetArticle(r.Context(), database, devID, id)
			if err != nil {
				writeError(w, http.StatusNotFound, "not_found", "article not found")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": a.ID, "sourceId": a.SourceID, "canonicalUrl": a.CanonicalURL, "title": a.Title, "summary": a.Summary, "author": a.Author, "publishedAt": a.PublishedAt, "isRead": a.IsRead,
			})
		})

		r.Post("/{id}/read", func(w http.ResponseWriter, r *http.Request) {
			devID := r.Context().Value(ctxDeviceID{}).(int64)
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid id")
				return
			}
			var body struct {
				IsRead bool `json:"isRead"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid json")
				return
			}
			if err := db.SetReadState(r.Context(), database, devID, id, body.IsRead); err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "toggle fail")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})
}
