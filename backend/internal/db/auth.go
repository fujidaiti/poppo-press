package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fujidaiti/poppo-press/backend/internal/auth"
)

// GetUserPasswordHash returns the stored PHC password hash for the given username.
func GetUserPasswordHash(ctx context.Context, database *sql.DB, username string) (string, error) {
	var phc string
	err := database.QueryRowContext(ctx, "SELECT password_hash FROM user WHERE username = ?", username).Scan(&phc)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return phc, err
}

// CreateOrUpdateDeviceToken stores the hash for the device token; creates the device row if needed.
// Returns device id.
func CreateOrUpdateDeviceToken(ctx context.Context, database *sql.DB, deviceName, tokenHash string) (int64, error) {
	// Try update existing device by name
	res, err := database.ExecContext(ctx,
		"UPDATE device SET token_hash = ?, last_seen_at = ? WHERE name = ?",
		tokenHash, time.Now().UTC().Format(time.RFC3339), deviceName,
	)
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		var id int64
		err = database.QueryRowContext(ctx, "SELECT id FROM device WHERE name = ?", deviceName).Scan(&id)
		return id, err
	}
	// Insert new device
	res, err = database.ExecContext(ctx,
		"INSERT INTO device(name, token_hash, last_seen_at, created_at) VALUES(?,?,?,?)",
		deviceName, tokenHash, time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// LookupDeviceIdByToken returns device id matching the token hash and updates last_seen_at.
func LookupDeviceIdByToken(ctx context.Context, database *sql.DB, token string) (int64, error) {
	hash := auth.HashToken(token)
	var id int64
	err := database.QueryRowContext(ctx, "SELECT id FROM device WHERE token_hash = ?", hash).Scan(&id)
	if err != nil {
		return 0, err
	}
	_, _ = database.ExecContext(ctx, "UPDATE device SET last_seen_at = ? WHERE id = ?", time.Now().UTC().Format(time.RFC3339), id)
	return id, nil
}

// RevokeDeviceToken clears token_hash for the given device id.
func RevokeDeviceToken(ctx context.Context, database *sql.DB, deviceID int64) error {
	_, err := database.ExecContext(ctx, "UPDATE device SET token_hash = NULL WHERE id = ?", deviceID)
	return err
}
