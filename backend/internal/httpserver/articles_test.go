package httpserver

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fujidaiti/poppo-press/backend/internal/testutil"
)

func TestArticlesListDetailAndReadToggle(t *testing.T) {
	db, cleanup := testutil.OpenTestDB(t, "admin-pass")
	defer cleanup()

	// seed a source and two articles
	res, err := db.Exec("INSERT INTO source(url, created_at) VALUES(?, ?)", "https://ex/feed", time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert source: %v", err)
	}
	srcID, _ := res.LastInsertId()
	mustExec(t, db, `INSERT INTO article(id, source_id, canonical_url, title, summary, published_at, created_at, canonical_id) VALUES(?,?,?,?,?,?,?,?)`, 101, srcID, "https://ex/a", "A", "sa", time.Now().Add(-time.Hour).UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339), "aid-101")
	mustExec(t, db, `INSERT INTO article(id, source_id, canonical_url, title, summary, published_at, created_at, canonical_id) VALUES(?,?,?,?,?,?,?,?)`, 102, srcID, "https://ex/b", "B", "sb", time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339), "aid-102")

	srv := New(db)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	// login
	body := map[string]string{"username": "admin", "password": "admin-pass", "deviceName": "dev"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	var lr struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&lr)
	_ = resp.Body.Close()

	// list articles
	req1, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/articles", nil)
	req1.Header.Set("Authorization", "Bearer "+lr.Token)
	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("list status: %d", resp1.StatusCode)
	}
	var list []struct {
		ID     int64 `json:"id"`
		IsRead bool  `json:"isRead"`
	}
	if err := json.NewDecoder(resp1.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	_ = resp1.Body.Close()
	if len(list) != 2 || list[0].IsRead || list[1].IsRead {
		t.Fatalf("expected unread list, got %+v", list)
	}

	// mark 101 read
	rb, _ := json.Marshal(map[string]bool{"isRead": true})
	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/articles/101/read", bytes.NewReader(rb))
	req2.Header.Set("Authorization", "Bearer "+lr.Token)
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("toggle: %v", err)
	}
	if resp2.StatusCode != http.StatusNoContent {
		t.Fatalf("toggle status: %d", resp2.StatusCode)
	}
	_ = resp2.Body.Close()

	// detail 101 should show isRead true
	req3, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/articles/101", nil)
	req3.Header.Set("Authorization", "Bearer "+lr.Token)
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("detail status: %d", resp3.StatusCode)
	}
	var ad struct {
		ID     int64 `json:"id"`
		IsRead bool  `json:"isRead"`
	}
	if err := json.NewDecoder(resp3.Body).Decode(&ad); err != nil {
		t.Fatalf("decode detail: %v", err)
	}
	_ = resp3.Body.Close()
	if !ad.IsRead {
		t.Fatalf("expected isRead true, got false")
	}

	// filter read
	req4, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/articles?readState=read", nil)
	req4.Header.Set("Authorization", "Bearer "+lr.Token)
	resp4, err := http.DefaultClient.Do(req4)
	if err != nil {
		t.Fatalf("list read: %v", err)
	}
	var listR []struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(resp4.Body).Decode(&listR)
	_ = resp4.Body.Close()
	if len(listR) != 1 || listR[0].ID != 101 {
		t.Fatalf("expected only 101 read, got %+v", listR)
	}

	// mark 101 unread
	rb2, _ := json.Marshal(map[string]bool{"isRead": false})
	req5, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/articles/101/read", bytes.NewReader(rb2))
	req5.Header.Set("Authorization", "Bearer "+lr.Token)
	req5.Header.Set("Content-Type", "application/json")
	resp5, err := http.DefaultClient.Do(req5)
	if err != nil {
		t.Fatalf("toggle2: %v", err)
	}
	if resp5.StatusCode != http.StatusNoContent {
		t.Fatalf("toggle2 status: %d", resp5.StatusCode)
	}
	_ = resp5.Body.Close()

	// filter unread should return both
	req6, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/articles?readState=unread", nil)
	req6.Header.Set("Authorization", "Bearer "+lr.Token)
	resp6, err := http.DefaultClient.Do(req6)
	if err != nil {
		t.Fatalf("list unread: %v", err)
	}
	var listU []struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(resp6.Body).Decode(&listU)
	_ = resp6.Body.Close()
	if len(listU) != 2 {
		t.Fatalf("expected 2 unread, got %d", len(listU))
	}
}

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec: %v", err)
	}
}
