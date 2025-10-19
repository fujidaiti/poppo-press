package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
