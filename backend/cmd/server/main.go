package main

import (
	"log"
	"net/http"

	"github.com/fujidaiti/poppo-press/backend/internal/config"
	"github.com/fujidaiti/poppo-press/backend/internal/db"
	"github.com/fujidaiti/poppo-press/backend/internal/httpserver"
)

func main() {
	cfg := config.Load()
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	srv := httpserver.New()
	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
