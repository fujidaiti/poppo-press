package db

import (
	"context"
	"database/sql"
)

type UpsertArticleParams struct {
	SourceID     int64
	CanonicalURL string
	Title        string
	Summary      string
	Content      string
	Author       string
	PublishedAt  string
	UpdatedAt    string
	CanonicalID  string
}

func UpsertArticleByCanonicalID(ctx context.Context, database *sql.DB, p UpsertArticleParams) error {
	// Upsert by canonical_id (unique when not null)
	if p.CanonicalID == "" {
		return nil
	}
	_, err := database.ExecContext(ctx, `
INSERT INTO article(source_id, canonical_url, title, summary, content, author, published_at, updated_at, canonical_id)
VALUES(?,?,?,?,?,?,?,?,?)
ON CONFLICT(canonical_id) DO UPDATE SET
  title=excluded.title,
  summary=excluded.summary,
  content=excluded.content,
  author=excluded.author,
  published_at=excluded.published_at,
  updated_at=excluded.updated_at
`, p.SourceID, p.CanonicalURL, p.Title, p.Summary, p.Content, p.Author, p.PublishedAt, p.UpdatedAt, p.CanonicalID)
	return err
}
