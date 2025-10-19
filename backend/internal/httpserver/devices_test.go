package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fujidaiti/poppo-press/backend/internal/testutil"
)

func TestDevicesListAndRevoke(t *testing.T) {
	db, cleanup := testutil.OpenTestDB(t, "admin-pass")
	defer cleanup()

	srv := New(db)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	// Login to create a device
	body := map[string]string{"username": "admin", "password": "admin-pass", "deviceName": "dev-1"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status: %d", resp.StatusCode)
	}
	var lr struct {
		Token    string `json:"token"`
		DeviceID int64  `json:"deviceId"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&lr)
	_ = resp.Body.Close()

	// List devices should include our device
	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/devices", nil)
	req2.Header.Set("Authorization", "Bearer "+lr.Token)
	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("list status: %d", resp2.StatusCode)
	}
	var list []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	_ = json.NewDecoder(resp2.Body).Decode(&list)
	_ = resp2.Body.Close()
	if len(list) == 0 {
		t.Fatalf("expected at least one device")
	}
	found := false
	for _, d := range list {
		if d.ID == lr.DeviceID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("device not in list: %+v", list)
	}

	// Revoke device
	req3, _ := http.NewRequest(http.MethodDelete, ts.URL+"/v1/devices/"+itoa(lr.DeviceID), nil)
	req3.Header.Set("Authorization", "Bearer "+lr.Token)
	resp3, err := http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if resp3.StatusCode != http.StatusNoContent {
		t.Fatalf("revoke status: %d", resp3.StatusCode)
	}
	_ = resp3.Body.Close()

	// Using the same token should now fail
	req4, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/protected/ping", nil)
	req4.Header.Set("Authorization", "Bearer "+lr.Token)
	resp4, err := http.DefaultClient.Do(req4)
	if err != nil {
		t.Fatalf("protected: %v", err)
	}
	if resp4.StatusCode != http.StatusUnauthorized {
		t.Fatalf("protected status: %d", resp4.StatusCode)
	}
	_ = resp4.Body.Close()
}

func itoa(v int64) string { return fmt.Sprintf("%d", v) }
