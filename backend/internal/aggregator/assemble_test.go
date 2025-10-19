package aggregator

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/fujidaiti/poppo-press/backend/internal/testutil"
)

func TestAssembleDailyEdition_IdempotentAndOrdered(t *testing.T) {
	db, cleanup := testutil.OpenTestDB(t, "admin-pass")
	defer cleanup()

	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2025, 10, 19, 8, 0, 0, 0, time.UTC)

	srcID := insertSource(t, db, "https://ex/feed")
	// three items: newest first should get position 1
	insertArticle(t, db, srcID, 1, now.Add(-2*time.Hour))
	insertArticle(t, db, srcID, 2, now.Add(-1*time.Hour))
	insertArticle(t, db, srcID, 3, now.Add(-30*time.Minute))

	if err := AssembleDailyEdition(context.Background(), db, loc, now); err != nil {
		t.Fatalf("assemble 1: %v", err)
	}

	edID := getEditionID(t, db, "2025-10-19")
	assertPositions(t, db, edID, []int64{3, 2, 1})

	// Re-run is idempotent
	if err := AssembleDailyEdition(context.Background(), db, loc, now); err != nil {
		t.Fatalf("assemble 2: %v", err)
	}
	edID2 := getEditionID(t, db, "2025-10-19")
	if edID2 != edID {
		t.Fatalf("edition id changed: %d vs %d", edID2, edID)
	}
	assertPositions(t, db, edID, []int64{3, 2, 1})

	// Old article outside 24h window should be excluded
	insertArticle(t, db, srcID, 4, now.Add(-25*time.Hour))
	if err := AssembleDailyEdition(context.Background(), db, loc, now); err != nil {
		t.Fatalf("assemble 3: %v", err)
	}
	assertPositions(t, db, edID, []int64{3, 2, 1})
}

func insertSource(t *testing.T, db *sql.DB, url string) int64 {
	t.Helper()
	res, err := db.Exec("INSERT INTO source(url, created_at) VALUES(?, ?)", url, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert source: %v", err)
	}
	id, _ := res.LastInsertId()
	return id
}

func insertArticle(t *testing.T, db *sql.DB, srcID int64, aid int64, published time.Time) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO article(id, source_id, canonical_url, title, published_at, created_at, canonical_id)
VALUES(?,?,?,?,?,?,?)`, aid, srcID, "https://ex/a", "t", published.UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339), aid)
	if err != nil {
		t.Fatalf("insert article: %v", err)
	}
}

func getEditionID(t *testing.T, db *sql.DB, localDate string) int64 {
	t.Helper()
	var id int64
	if err := db.QueryRow("SELECT id FROM edition WHERE local_date = ?", localDate).Scan(&id); err != nil {
		t.Fatalf("get edition: %v", err)
	}
	return id
}

func assertPositions(t *testing.T, db *sql.DB, edID int64, want []int64) {
	t.Helper()
	rows, err := db.Query("SELECT article_id FROM edition_article WHERE edition_id = ? ORDER BY position", edID)
	if err != nil {
		t.Fatalf("query pos: %v", err)
	}
	defer rows.Close()
	var got []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("scan: %v", err)
		}
		got = append(got, id)
	}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: %v vs %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("pos %d: got %d want %d", i, got[i], want[i])
		}
	}
}
