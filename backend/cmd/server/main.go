package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fujidaiti/poppo-press/backend/internal/config"
	"github.com/fujidaiti/poppo-press/backend/internal/db"
	"github.com/fujidaiti/poppo-press/backend/internal/httpserver"
)

// main loads configuration, initializes the database (migrate and seed), and
// starts the HTTP server with health and version endpoints.
func main() {
	cfg := config.Load()
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.Migrate(context.Background(), database); err != nil {
		log.Fatal(err)
	}
	if err := db.SeedAdminIfEmpty(context.Background(), database); err != nil {
		log.Fatal(err)
	}

	srv := httpserver.New(database)
	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
