// Package db provides database initialization, migrations, and helpers for the
// SQLite-backed storage layer.
package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/fujidaiti/poppo-press/backend/internal/auth"
)

// SeedAdminIfEmpty ensures there is at least one user in the database. If no
// users exist, it creates an admin user using PP_ADMIN_USER (default "admin")
// and PP_ADMIN_PASS (required). The password is hashed using argon2id.
func SeedAdminIfEmpty(ctx context.Context, database *sql.DB) error {
	var count int
	if err := database.QueryRowContext(ctx, "SELECT COUNT(1) FROM user").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	username := os.Getenv("PP_ADMIN_USER")
	if username == "" {
		username = "admin"
	}
	pass := os.Getenv("PP_ADMIN_PASS")
	if pass == "" {
		return errors.New("PP_ADMIN_PASS is required for first-run seeding")
	}

	phc, err := auth.HashPassword(pass)
	if err != nil {
		return err
	}
	_, err = database.ExecContext(ctx,
		"INSERT INTO user(username, password_hash, created_at) VALUES(?,?,?)",
		username, phc, time.Now().UTC().Format(time.RFC3339),
	)
	return err
}
