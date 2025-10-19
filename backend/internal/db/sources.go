package db

import (
	"context"
	"database/sql"
)

type SourceRow struct {
	ID        int64
	URL       string
	Title     string
	CreatedAt string
}

type SourceFetchRow struct {
	ID           int64
	URL          string
	ETag         string
	LastModified string
}

func CreateSource(ctx context.Context, database *sql.DB, url, title, etag, lastModified string) (int64, error) {
	res, err := database.ExecContext(ctx,
		"INSERT INTO source(url, title, etag, last_modified) VALUES(?,?,?,?)",
		url, title, etag, lastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ListSources(ctx context.Context, database *sql.DB) ([]SourceRow, error) {
	rows, err := database.QueryContext(ctx, "SELECT id, url, IFNULL(title, ''), created_at FROM source ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SourceRow
	for rows.Next() {
		var r SourceRow
		if err := rows.Scan(&r.ID, &r.URL, &r.Title, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func ListSourcesForFetch(ctx context.Context, database *sql.DB) ([]SourceFetchRow, error) {
	rows, err := database.QueryContext(ctx, "SELECT id, url, IFNULL(etag, ''), IFNULL(last_modified, '') FROM source ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SourceFetchRow
	for rows.Next() {
		var r SourceFetchRow
		if err := rows.Scan(&r.ID, &r.URL, &r.ETag, &r.LastModified); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func UpdateSourceHeaders(ctx context.Context, database *sql.DB, id int64, etag, lastModified string) error {
	_, err := database.ExecContext(ctx, "UPDATE source SET etag = ?, last_modified = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", etag, lastModified, id)
	return err
}

func DeleteSource(ctx context.Context, database *sql.DB, id int64) (bool, error) {
	res, err := database.ExecContext(ctx, "DELETE FROM source WHERE id = ?", id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}
