package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
)

// migrationsFS embeds SQL migrations stored under migrations/.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate applies any pending embedded SQL migrations in lexicographic order.
// It is idempotent, tracking the latest applied version in schema_migrations.
func Migrate(ctx context.Context, database *sql.DB) error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER NOT NULL)"); err != nil {
		return err
	}

	var current int
	_ = tx.QueryRowContext(ctx, "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&current)

	for _, name := range files {
		var ver int
		if _, err := fmt.Sscanf(name, "%d_", &ver); err != nil {
			continue
		}
		if ver <= current {
			continue
		}
		b, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, string(b)); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations(version) VALUES (?)", ver); err != nil {
			return err
		}
		current = ver
	}

	return tx.Commit()
}
