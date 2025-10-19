package db

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	dsn := path + "?_busy_timeout=5000&_fk=1"
	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(1)
	database.SetConnMaxIdleTime(5 * time.Minute)
	database.SetConnMaxLifetime(0)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := pragma(ctx, database, "PRAGMA journal_mode=WAL"); err != nil {
		_ = database.Close()
		return nil, err
	}
	if err := pragma(ctx, database, "PRAGMA synchronous=NORMAL"); err != nil {
		_ = database.Close()
		return nil, err
	}
	return database, nil
}

func pragma(ctx context.Context, db *sql.DB, stmt string) error {
	_, err := db.ExecContext(ctx, stmt)
	return err
}
