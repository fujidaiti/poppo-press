package fetcher

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/fujidaiti/poppo-press/backend/internal/testutil"
)

func TestFetchAllSources_ConditionalGET_ParseAndDedupe(t *testing.T) {
	database, cleanup := testutil.OpenTestDB(t, "admin-pass")
	defer cleanup()

	// Start with a dynamic feed server; allow swapping handler behavior
	etag := "W/\"v1\""
	lastMod := "Mon, 06 Sep 2021 00:00:00 GMT"
	items := []string{
		`<item><guid>1</guid><title>A</title><link>https://ex/a</link><pubDate>Mon, 06 Sep 2021 00:00:00 GMT</pubDate></item>`,
		`<item><guid>2</guid><title>B</title><link>https://ex/b</link><pubDate>Mon, 06 Sep 2021 01:00:00 GMT</pubDate></item>`,
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Support If-None-Match/If-Modified-Since naive check
		if r.Header.Get("If-None-Match") == etag || r.Header.Get("If-Modified-Since") == lastMod {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("ETag", etag)
		w.Header().Set("Last-Modified", lastMod)
		w.Header().Set("Content-Type", "application/rss+xml")
		fmt.Fprintf(w, `<?xml version="1.0"?><rss version="2.0"><channel><title>T</title>%s</channel></rss>`, join(items))
	}
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	// Insert one source pointing to the feed
	srcID := insertSource(t, database, srv.URL)

	// First run: should fetch, parse 2 items, update headers
	if err := FetchAllSources(context.Background(), database, &http.Client{Timeout: 5 * time.Second}); err != nil {
		t.Fatalf("fetch 1: %v", err)
	}
	assertArticleCount(t, database, srcID, 2)
	assertSourceHeaders(t, database, srcID, etag, lastMod)

	// Second run with 304: no new items
	if err := FetchAllSources(context.Background(), database, &http.Client{Timeout: 5 * time.Second}); err != nil {
		t.Fatalf("fetch 2: %v", err)
	}
	assertArticleCount(t, database, srcID, 2)

	// Change feed: new item and new ETag
	etag = "W/\"v2\""
	items = append(items, `<item><guid>3</guid><title>C</title><link>https://ex/c</link><pubDate>Mon, 06 Sep 2021 02:00:00 GMT</pubDate></item>`)

	if err := FetchAllSources(context.Background(), database, &http.Client{Timeout: 5 * time.Second}); err != nil {
		t.Fatalf("fetch 3: %v", err)
	}
	assertArticleCount(t, database, srcID, 3)
}

func join(items []string) string {
	s := ""
	for _, it := range items {
		s += it
	}
	return s
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

func assertArticleCount(t *testing.T, db *sql.DB, sourceID int64, want int) {
	t.Helper()
	var got int
	if err := db.QueryRow("SELECT COUNT(1) FROM article WHERE source_id = ?", sourceID).Scan(&got); err != nil {
		t.Fatalf("count: %v", err)
	}
	if got != want {
		t.Fatalf("count mismatch: got %d want %d", got, want)
	}
}

func assertSourceHeaders(t *testing.T, db *sql.DB, id int64, etag, lastMod string) {
	t.Helper()
	var gotE, gotL string
	if err := db.QueryRow("SELECT etag, last_modified FROM source WHERE id = ?", id).Scan(&gotE, &gotL); err != nil {
		t.Fatalf("headers: %v", err)
	}
	if gotE != etag || gotL != lastMod {
		t.Fatalf("headers mismatch: %q %q vs %q %q", gotE, gotL, etag, lastMod)
	}
}
