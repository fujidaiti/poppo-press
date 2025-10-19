package db

import (
	"context"
	"database/sql"
)

type BookmarkRow struct {
	ArticleID int64
}

func ListBookmarks(ctx context.Context, database *sql.DB) ([]BookmarkRow, error) {
	rows, err := database.QueryContext(ctx, `SELECT article_id FROM bookmark ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BookmarkRow
	for rows.Next() {
		var r BookmarkRow
		if err := rows.Scan(&r.ArticleID); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func AddBookmark(ctx context.Context, database *sql.DB, articleID int64) error {
	_, err := database.ExecContext(ctx, `INSERT OR IGNORE INTO bookmark(article_id) VALUES(?)`, articleID)
	return err
}

func RemoveBookmark(ctx context.Context, database *sql.DB, articleID int64) error {
	_, err := database.ExecContext(ctx, `DELETE FROM bookmark WHERE article_id = ?`, articleID)
	return err
}
