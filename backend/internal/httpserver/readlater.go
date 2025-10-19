package httpserver

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fujidaiti/poppo-press/backend/internal/db"
)

func registerReadLaterRoutes(database *sql.DB, r chi.Router) {
	r.With(authMiddleware(database)).Route("/read-later", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			rows, err := db.ListBookmarks(r.Context(), database)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "list fail")
				return
			}
			type out struct {
				ID int64 `json:"id"`
			}
			outList := make([]out, 0, len(rows))
			for _, b := range rows {
				outList = append(outList, out{ID: b.ArticleID})
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(outList)
		})
		r.Post("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid id")
				return
			}
			if err := db.AddBookmark(r.Context(), database, id); err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "add fail")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid id")
				return
			}
			if err := db.RemoveBookmark(r.Context(), database, id); err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "delete fail")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})
}
