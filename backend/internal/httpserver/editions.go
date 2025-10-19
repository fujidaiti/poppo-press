package httpserver

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func registerEditionRoutes(database *sql.DB, r chi.Router) {
	r.With(authMiddleware(database)).Route("/editions", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			rows, err := database.QueryContext(r.Context(), `
SELECT e.id, e.local_date, e.published_at, COUNT(ea.article_id) as cnt
FROM edition e LEFT JOIN edition_article ea ON e.id = ea.edition_id
GROUP BY e.id, e.local_date, e.published_at
ORDER BY e.id DESC`)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "list fail")
				return
			}
			type out struct {
				ID           int64   `json:"id"`
				LocalDate    string  `json:"localDate"`
				PublishedAt  *string `json:"publishedAt"`
				ArticleCount int     `json:"articleCount"`
			}
			var list []out
			for rows.Next() {
				var o out
				var pub sql.NullString
				if err := rows.Scan(&o.ID, &o.LocalDate, &pub, &o.ArticleCount); err != nil {
					writeError(w, http.StatusInternalServerError, "internal", "scan fail")
					return
				}
				if pub.Valid {
					o.PublishedAt = &pub.String
				}
				list = append(list, o)
			}
			_ = rows.Err()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(list)
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			idStr := chi.URLParam(r, "id")
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid id")
				return
			}
			var localDate string
			var publishedAt *string
			var pub sql.NullString
			if err := database.QueryRowContext(r.Context(), `SELECT local_date, published_at FROM edition WHERE id = ?`, id).Scan(&localDate, &pub); err != nil {
				writeError(w, http.StatusNotFound, "not_found", "edition not found")
				return
			}
			if pub.Valid {
				v := pub.String
				publishedAt = &v
			}
			rows, err := database.QueryContext(r.Context(), `
SELECT a.id, a.source_id, a.canonical_url, a.title, a.summary, a.author, a.published_at
FROM edition_article ea JOIN article a ON ea.article_id = a.id
WHERE ea.edition_id = ? ORDER BY ea.position`, id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "query fail")
				return
			}
			type art struct {
				ID           int64  `json:"id"`
				SourceID     int64  `json:"sourceId"`
				CanonicalURL string `json:"canonicalUrl"`
				Title        string `json:"title"`
				Summary      string `json:"summary"`
				Author       string `json:"author"`
				PublishedAt  string `json:"publishedAt"`
			}
			var arts []art
			for rows.Next() {
				var a art
				if err := rows.Scan(&a.ID, &a.SourceID, &a.CanonicalURL, &a.Title, &a.Summary, &a.Author, &a.PublishedAt); err != nil {
					writeError(w, http.StatusInternalServerError, "internal", "scan fail")
					return
				}
				arts = append(arts, a)
			}
			_ = rows.Err()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": id, "localDate": localDate, "publishedAt": publishedAt, "articles": arts,
			})
		})
	})
}
