package httpserver

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fujidaiti/poppo-press/backend/internal/db"
)

func registerDeviceRoutes(database *sql.DB, r chi.Router) {
	r.With(authMiddleware(database)).Route("/devices", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			rows, err := db.ListDevices(r.Context(), database)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "list fail")
				return
			}
			type out struct {
				ID         int64  `json:"id"`
				Name       string `json:"name"`
				LastSeenAt string `json:"lastSeenAt"`
				CreatedAt  string `json:"createdAt"`
			}
			resp := make([]out, 0, len(rows))
			for _, d := range rows {
				resp = append(resp, out{ID: d.ID, Name: d.Name, LastSeenAt: d.LastSeenAt, CreatedAt: d.CreatedAt})
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		})
		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				writeError(w, http.StatusBadRequest, "bad_request", "invalid id")
				return
			}
			if err := db.RevokeDeviceToken(r.Context(), database, id); err != nil {
				writeError(w, http.StatusInternalServerError, "internal", "revoke fail")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		})
	})
}
