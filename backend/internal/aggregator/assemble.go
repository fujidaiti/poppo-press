package aggregator

import (
	"context"
	"database/sql"
	"time"
)

func AssembleDailyEdition(ctx context.Context, database *sql.DB, tz *time.Location, now time.Time) error {
	localDate := now.In(tz).Format("2006-01-02")
	windowStart := now.Add(-24 * time.Hour).UTC().Format(time.RFC3339)
	windowEnd := now.UTC().Format(time.RFC3339)

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Upsert edition row by local_date
	if _, err := tx.ExecContext(ctx, `
INSERT INTO edition(local_date, published_at, created_at)
VALUES(?,?,COALESCE(created_at, CURRENT_TIMESTAMP))
ON CONFLICT(local_date) DO UPDATE SET published_at=excluded.published_at
`, localDate, now.UTC().Format(time.RFC3339)); err != nil {
		return err
	}
	var editionID int64
	if err := tx.QueryRowContext(ctx, `SELECT id FROM edition WHERE local_date = ?`, localDate).Scan(&editionID); err != nil {
		return err
	}

	// Clear previous links
	if _, err := tx.ExecContext(ctx, `DELETE FROM edition_article WHERE edition_id = ?`, editionID); err != nil {
		return err
	}

	// Select last 24h articles ordered newest first
	rows, err := tx.QueryContext(ctx, `
SELECT id FROM article
WHERE published_at >= ? AND published_at <= ?
ORDER BY published_at DESC
`, windowStart, windowEnd)
	if err != nil {
		return err
	}
	defer rows.Close()

	position := 1
	for rows.Next() {
		var articleID int64
		if err := rows.Scan(&articleID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO edition_article(edition_id, article_id, position)
VALUES(?,?,?)
`, editionID, articleID, position); err != nil {
			return err
		}
		position++
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return tx.Commit()
}
