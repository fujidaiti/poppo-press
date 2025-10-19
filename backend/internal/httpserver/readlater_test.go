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

func TestReadLaterFlow(t *testing.T) {
	db, cleanup := testutil.OpenTestDB(t, "admin-pass")
	defer cleanup()

	// seed a source and two articles
	res, err := db.Exec("INSERT INTO source(url, created_at) VALUES(?, ?)", "https://ex/feed", time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("insert source: %v", err)
	}
	srcID, _ := res.LastInsertId()
	mustExecRL(t, db, `INSERT INTO article(id, source_id, canonical_url, title, summary, published_at, created_at, canonical_id) VALUES(?,?,?,?,?,?,?,?)`, 201, srcID, "https://ex/a", "A", "sa", time.Now().Add(-time.Hour).UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339), "aid-201")
	mustExecRL(t, db, `INSERT INTO article(id, source_id, canonical_url, title, summary, published_at, created_at, canonical_id) VALUES(?,?,?,?,?,?,?,?)`, 202, srcID, "https://ex/b", "B", "sb", time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339), "aid-202")

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

	// initially empty
	req0, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/read-later", nil)
	req0.Header.Set("Authorization", "Bearer "+lr.Token)
	resp0, err := http.DefaultClient.Do(req0)
	if err != nil {
		t.Fatalf("get empty: %v", err)
	}
	if resp0.StatusCode != http.StatusOK {
		t.Fatalf("status empty: %d", resp0.StatusCode)
	}
	var empty []any
	_ = json.NewDecoder(resp0.Body).Decode(&empty)
	_ = resp0.Body.Close()
	if len(empty) != 0 {
		t.Fatalf("expected empty, got %d", len(empty))
	}

	// add 201
	req1, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/read-later/201", nil)
	req1.Header.Set("Authorization", "Bearer "+lr.Token)
	resp1, err := http.DefaultClient.Do(req1)
	if err != nil {
		t.Fatalf("add 201: %v", err)
	}
	if resp1.StatusCode != http.StatusNoContent {
		t.Fatalf("add status: %d", resp1.StatusCode)
	}
	_ = resp1.Body.Close()

	// list should have 1
	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/read-later", nil)
	req2.Header.Set("Authorization", "Bearer "+lr.Token)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("list 1: %v", err)
	}
	var list1 []struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(resp2.Body).Decode(&list1)
	_ = resp2.Body.Close()
	if len(list1) != 1 || list1[0].ID != 201 {
		t.Fatalf("unexpected list1: %+v", list1)
	}

	// add 202
	req3, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/read-later/202", nil)
	req3.Header.Set("Authorization", "Bearer "+lr.Token)
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("add 202: %v", err)
	}
	if resp3.StatusCode != http.StatusNoContent {
		t.Fatalf("add 202 status: %d", resp3.StatusCode)
	}
	_ = resp3.Body.Close()

	// list should have 2
	req4, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/read-later", nil)
	req4.Header.Set("Authorization", "Bearer "+lr.Token)
	resp4, err := http.DefaultClient.Do(req4)
	if err != nil {
		t.Fatalf("list 2: %v", err)
	}
	var list2 []struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(resp4.Body).Decode(&list2)
	_ = resp4.Body.Close()
	if len(list2) != 2 {
		t.Fatalf("unexpected list2 len: %d", len(list2))
	}

	// delete 201
	req5, _ := http.NewRequest(http.MethodDelete, ts.URL+"/v1/read-later/201", nil)
	req5.Header.Set("Authorization", "Bearer "+lr.Token)
	resp5, err := http.DefaultClient.Do(req5)
	if err != nil {
		t.Fatalf("delete 201: %v", err)
	}
	if resp5.StatusCode != http.StatusNoContent {
		t.Fatalf("delete status: %d", resp5.StatusCode)
	}
	_ = resp5.Body.Close()

	// list should have only 202
	req6, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/read-later", nil)
	req6.Header.Set("Authorization", "Bearer "+lr.Token)
	resp6, err := http.DefaultClient.Do(req6)
	if err != nil {
		t.Fatalf("list 3: %v", err)
	}
	var list3 []struct {
		ID int64 `json:"id"`
	}
	_ = json.NewDecoder(resp6.Body).Decode(&list3)
	_ = resp6.Body.Close()
	if len(list3) != 1 || list3[0].ID != 202 {
		t.Fatalf("unexpected list3: %+v", list3)
	}
}

func mustExecRL(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec: %v", err)
	}
}
