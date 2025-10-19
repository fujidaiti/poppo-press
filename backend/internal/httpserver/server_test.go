package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"database/sql"

	"github.com/fujidaiti/poppo-press/backend/internal/testutil"
)

func TestLoginLogoutAndProtectedRoute(t *testing.T) {
	db, cleanup := testutil.OpenTestDB(t, "admin-pass")
	defer cleanup()

	srv := New(db)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	// Login
	body := map[string]string{"username": "admin", "password": "admin-pass", "deviceName": "test-device"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("login request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status: %d", resp.StatusCode)
	}
	var loginResp struct {
		Token    string `json:"token"`
		DeviceID int64  `json:"deviceId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	_ = resp.Body.Close()
	if loginResp.Token == "" || loginResp.DeviceID == 0 {
		t.Fatalf("invalid login response: %+v", loginResp)
	}

	// Protected ping should work
	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/protected/ping", nil)
	req2.Header.Set("Authorization", "Bearer "+loginResp.Token)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("protected request: %v", err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("protected status: %d", resp2.StatusCode)
	}
	_ = resp2.Body.Close()

	// Logout
	req3, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/auth/logout", nil)
	req3.Header.Set("Authorization", "Bearer "+loginResp.Token)
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("logout request: %v", err)
	}
	if resp3.StatusCode != http.StatusNoContent {
		t.Fatalf("logout status: %d", resp3.StatusCode)
	}
	_ = resp3.Body.Close()

	// Protected should now fail
	req4, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/protected/ping", nil)
	req4.Header.Set("Authorization", "Bearer "+loginResp.Token)
	resp4, err := http.DefaultClient.Do(req4)
	if err != nil {
		t.Fatalf("protected request 2: %v", err)
	}
	if resp4.StatusCode != http.StatusUnauthorized {
		t.Fatalf("protected status after logout: %d", resp4.StatusCode)
	}
	_ = resp4.Body.Close()
}

func TestSourcesAPI(t *testing.T) {
	db, cleanup := testutil.OpenTestDB(t, "admin-pass")
	defer cleanup()

	srv := New(db)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	// login to get token
	loginBody := map[string]string{"username": "admin", "password": "admin-pass", "deviceName": "dev"}
	lb, _ := json.Marshal(loginBody)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/auth/login", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("login request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status: %d", resp.StatusCode)
	}
	var lr struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	_ = resp.Body.Close()

	// invalid URL should 400
	bad := map[string]string{"url": "not-a-url"}
	bb, _ := json.Marshal(bad)
	req2, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/sources", bytes.NewReader(bb))
	req2.Header.Set("Authorization", "Bearer "+lr.Token)
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("post bad url: %v", err)
	}
	if resp2.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad url status: %d", resp2.StatusCode)
	}
	_ = resp2.Body.Close()

	// spin up a simple RSS feed server with ETag/Last-Modified
	feedXML := "<?xml version=\"1.0\"?><rss version=\"2.0\"><channel><title>Test Feed</title></channel></rss>"
	etag := "W/\"abc123\""
	lastMod := "Mon, 06 Sep 2021 00:00:00 GMT"
	feed := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", etag)
		w.Header().Set("Last-Modified", lastMod)
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(feedXML))
	}))
	defer feed.Close()

	// create source
	good := map[string]string{"url": feed.URL}
	gb, _ := json.Marshal(good)
	req3, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/sources", bytes.NewReader(gb))
	req3.Header.Set("Authorization", "Bearer "+lr.Token)
	req3.Header.Set("Content-Type", "application/json")
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("post source: %v", err)
	}
	if resp3.StatusCode != http.StatusCreated {
		t.Fatalf("create status: %d", resp3.StatusCode)
	}
	var cr struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(resp3.Body).Decode(&cr); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	_ = resp3.Body.Close()
	if cr.ID == 0 {
		t.Fatalf("invalid created id: %d", cr.ID)
	}

	// GET list should contain it with title
	req4, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/sources", nil)
	req4.Header.Set("Authorization", "Bearer "+lr.Token)
	resp4, err := http.DefaultClient.Do(req4)
	if err != nil {
		t.Fatalf("get sources: %v", err)
	}
	if resp4.StatusCode != http.StatusOK {
		t.Fatalf("list status: %d", resp4.StatusCode)
	}
	var list []struct {
		ID        int64  `json:"id"`
		URL       string `json:"url"`
		Title     string `json:"title"`
		CreatedAt string `json:"createdAt"`
	}
	if err := json.NewDecoder(resp4.Body).Decode(&list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	_ = resp4.Body.Close()
	if len(list) != 1 || list[0].ID != cr.ID || list[0].URL != feed.URL || list[0].Title != "Test Feed" {
		t.Fatalf("unexpected list: %+v", list)
	}

	// ensure etag/last_modified persisted
	var gotETag, gotLM string
	if err := querySourceHeaders(t, db, cr.ID, &gotETag, &gotLM); err != nil {
		t.Fatalf("query headers: %v", err)
	}
	if gotETag != etag || gotLM != lastMod {
		t.Fatalf("headers mismatch: %q %q", gotETag, gotLM)
	}

	// DELETE
	req5, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf(ts.URL+"/v1/sources/%d", cr.ID), nil)
	req5.Header.Set("Authorization", "Bearer "+lr.Token)
	resp5, err := http.DefaultClient.Do(req5)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if resp5.StatusCode != http.StatusNoContent {
		t.Fatalf("delete status: %d", resp5.StatusCode)
	}
	_ = resp5.Body.Close()

	// list empty
	req6, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/sources", nil)
	req6.Header.Set("Authorization", "Bearer "+lr.Token)
	resp6, err := http.DefaultClient.Do(req6)
	if err != nil {
		t.Fatalf("get sources 2: %v", err)
	}
	if resp6.StatusCode != http.StatusOK {
		t.Fatalf("list 2 status: %d", resp6.StatusCode)
	}
	var list2 []any
	if err := json.NewDecoder(resp6.Body).Decode(&list2); err != nil {
		t.Fatalf("decode list2: %v", err)
	}
	_ = resp6.Body.Close()
	if len(list2) != 0 {
		t.Fatalf("expected empty, got %d", len(list2))
	}
}

func querySourceHeaders(t *testing.T, database *sql.DB, id int64, etag, lastMod *string) error {
	t.Helper()
	ctx := context.Background()
	return database.QueryRowContext(ctx, "SELECT etag, last_modified FROM source WHERE id = ?", id).Scan(etag, lastMod)
}
