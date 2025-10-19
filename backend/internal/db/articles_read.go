package db

import (
	"context"
	"database/sql"
)

type ArticleListRow struct {
	ID           int64
	SourceID     int64
	CanonicalURL string
	Title        string
	Summary      string
	Author       string
	PublishedAt  string
	IsRead       bool
}

func ListArticles(ctx context.Context, database *sql.DB, deviceID int64, readState string) ([]ArticleListRow, error) {
	// readState: "read" | "unread" | "all"
	where := ""
	if readState == "read" {
		where = "WHERE rs.is_read = 1"
	} else if readState == "unread" {
		where = "WHERE rs.is_read IS NULL OR rs.is_read = 0"
	}
	q := `
SELECT a.id, a.source_id, a.canonical_url, a.title, a.summary, a.author, a.published_at,
       COALESCE(rs.is_read, 0) as is_read
FROM article a
LEFT JOIN read_state rs ON rs.article_id = a.id AND rs.device_id = ?
` + where + `
ORDER BY a.published_at DESC`
	rows, err := database.QueryContext(ctx, q, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ArticleListRow
	for rows.Next() {
		var r ArticleListRow
		var isReadInt int
		if err := rows.Scan(&r.ID, &r.SourceID, &r.CanonicalURL, &r.Title, &r.Summary, &r.Author, &r.PublishedAt, &isReadInt); err != nil {
			return nil, err
		}
		r.IsRead = isReadInt == 1
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func GetArticle(ctx context.Context, database *sql.DB, deviceID, id int64) (ArticleListRow, error) {
	var r ArticleListRow
	var isReadInt int
	err := database.QueryRowContext(ctx, `
SELECT a.id, a.source_id, a.canonical_url, a.title, a.summary, a.author, a.published_at,
       COALESCE(rs.is_read, 0) as is_read
FROM article a
LEFT JOIN read_state rs ON rs.article_id = a.id AND rs.device_id = ?
WHERE a.id = ?
`, deviceID, id).Scan(&r.ID, &r.SourceID, &r.CanonicalURL, &r.Title, &r.Summary, &r.Author, &r.PublishedAt, &isReadInt)
	if err != nil {
		return ArticleListRow{}, err
	}
	r.IsRead = isReadInt == 1
	return r, nil
}

func SetReadState(ctx context.Context, database *sql.DB, deviceID, articleID int64, isRead bool) error {
	// upsert into read_state
	v := 0
	if isRead {
		v = 1
	}
	_, err := database.ExecContext(ctx, `
INSERT INTO read_state(article_id, device_id, is_read, updated_at)
VALUES(?,?,?,CURRENT_TIMESTAMP)
ON CONFLICT(article_id, device_id) DO UPDATE SET is_read=excluded.is_read, updated_at=excluded.updated_at
`, articleID, deviceID, v)
	return err
}
