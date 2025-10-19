package fetcher

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/fujidaiti/poppo-press/backend/internal/db"
)

func FetchAllSources(ctx context.Context, database *sql.DB, client *http.Client) error {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	sources, err := db.ListSourcesForFetch(ctx, database)
	if err != nil {
		return err
	}
	for _, s := range sources {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
		if err != nil {
			continue
		}
		if s.ETag != "" {
			req.Header.Set("If-None-Match", s.ETag)
		}
		if s.LastModified != "" {
			req.Header.Set("If-Modified-Since", s.LastModified)
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusNotModified {
			_ = resp.Body.Close()
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			_ = resp.Body.Close()
			continue
		}
		etag := resp.Header.Get("ETag")
		lastMod := resp.Header.Get("Last-Modified")
		parser := gofeed.NewParser()
		feed, err := parser.Parse(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			continue
		}
		for _, item := range feed.Items {
			published := time.Now().UTC()
			if item.PublishedParsed != nil {
				published = *item.PublishedParsed
			} else if item.UpdatedParsed != nil {
				published = *item.UpdatedParsed
			}
			author := ""
			if item.Author != nil {
				author = item.Author.Name
			}
			canonicalID := computeCanonicalID(item.GUID, item.Link, item.Title, published)
			_ = db.UpsertArticleByCanonicalID(ctx, database, db.UpsertArticleParams{
				SourceID:     s.ID,
				CanonicalURL: item.Link,
				Title:        item.Title,
				Summary:      item.Description,
				Content:      "",
				Author:       author,
				PublishedAt:  published.UTC().Format(time.RFC3339),
				UpdatedAt:    time.Now().UTC().Format(time.RFC3339),
				CanonicalID:  canonicalID,
			})
		}
		_ = db.UpdateSourceHeaders(ctx, database, s.ID, etag, lastMod)
	}
	return nil
}

func computeCanonicalID(guid, link, title string, published time.Time) string {
	if guid != "" {
		return guid
	}
	h := sha1.Sum([]byte(link + "|" + title + "|" + published.UTC().Format(time.RFC3339)))
	return hex.EncodeToString(h[:])
}
