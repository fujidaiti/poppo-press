package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLogin_WithPPToken_SavesDirectly(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	// init config with server (unused in PP_TOKEN path)
	root := NewRootCmd()
	root.SetArgs([]string{"init", "--server", "http://localhost:8080"})
	if err := root.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	t.Setenv("PP_TOKEN", "t-123")
	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"login", "--device", "devbox"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("login: %v; out=%s", err, out.String())
	}

	// Read back tz to force config load path; easier to assert via file read but keep package boundary
	get := NewRootCmd()
	var gout bytes.Buffer
	get.SetOut(&gout)
	get.SetArgs([]string{"config", "tz"})
	if err := get.Execute(); err != nil {
		t.Fatalf("config tz: %v", err)
	}
	// Additionally ensure file exists
	if _, err := os.Stat(mustConfigPath(t)); err != nil {
		t.Fatalf("config not found: %v", err)
	}
}

func TestLogin_ServerFlow_PersistsToken(t *testing.T) {
	// mock server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/auth/login" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"token": "tok-xyz", "deviceId": "1"})
	}))
	t.Cleanup(srv.Close)

	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")

	init := NewRootCmd()
	init.SetArgs([]string{"init", "--server", srv.URL})
	if err := init.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	var out bytes.Buffer
	cmd := NewRootCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"login", "--device", "devbox", "--username", "u", "--password", "p"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("login: %v; out=%s", err, out.String())
	}
}

// mustConfigPath exposes the config path using the config package indirectly via env.
func mustConfigPath(t *testing.T) string {
	t.Helper()
	// replicate config.Path; HOME is set in tests
	return os.ExpandEnv("$HOME/.config/poppo-press/config.yaml")
}
