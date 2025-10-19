package db

import (
	"context"
	"database/sql"
)

type DeviceRow struct {
	ID         int64
	Name       string
	LastSeenAt string
	CreatedAt  string
}

func ListDevices(ctx context.Context, database *sql.DB) ([]DeviceRow, error) {
	rows, err := database.QueryContext(ctx, `SELECT id, name, IFNULL(last_seen_at, ''), created_at FROM device ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DeviceRow
	for rows.Next() {
		var r DeviceRow
		if err := rows.Scan(&r.ID, &r.Name, &r.LastSeenAt, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
