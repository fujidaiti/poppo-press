package testutil

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fujidaiti/poppo-press/backend/internal/db"
)

// OpenTestDB opens a temporary SQLite database, runs migrations, seeds an admin
// user with the given password, and returns the DB and a cleanup function.
func OpenTestDB(t *testing.T, adminPassword string) (*sql.DB, func()) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	database, err := db.Open(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.Migrate(ctx, database); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	prev := os.Getenv("PP_ADMIN_PASS")
	_ = os.Setenv("PP_ADMIN_PASS", adminPassword)
	if err := db.SeedAdminIfEmpty(ctx, database); err != nil {
		t.Fatalf("seed admin: %v", err)
	}
	cleanup := func() {
		_ = os.Setenv("PP_ADMIN_PASS", prev)
		_ = database.Close()
	}
	return database, cleanup
}
