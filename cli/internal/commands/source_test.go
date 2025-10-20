package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSource_Add_List_Rm(t *testing.T) {
	// simple mock for /v1/sources and /v1/sources/{id}
	var createdId int64 = 42
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/sources":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": createdId})
			return
		case r.Method == http.MethodGet && r.URL.Path == "/v1/sources":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": createdId, "title": "Example", "url": "https://e/x"}})
			return
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/sources/42":
			w.WriteHeader(http.StatusNoContent)
			return
		}
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	t.Cleanup(srv.Close)

	// init config with server and token
	root := NewRootCmd()
	root.SetArgs([]string{"init", "--server", srv.URL})
	if err := root.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	// save token directly via PP_TOKEN path on login
	t.Setenv("PP_TOKEN", "tok")
	lg := NewRootCmd()
	lg.SetArgs([]string{"login", "--device", "dev"})
	if err := lg.Execute(); err != nil {
		t.Fatalf("login: %v", err)
	}

	// add
	var out bytes.Buffer
	add := NewRootCmd()
	add.SetOut(&out)
	add.SetArgs([]string{"source", "add", "https://e/x"})
	if err := add.Execute(); err != nil {
		t.Fatalf("add: %v; out=%s", err, out.String())
	}
	if got := out.String(); got == "" {
		t.Fatalf("expected output with created id, got empty")
	}

	// list
	out.Reset()
	ls := NewRootCmd()
	ls.SetOut(&out)
	ls.SetArgs([]string{"source", "list"})
	if err := ls.Execute(); err != nil {
		t.Fatalf("list: %v; out=%s", err, out.String())
	}
	if got := out.String(); got == "" {
		t.Fatalf("expected non-empty list output")
	}

	// rm
	rm := NewRootCmd()
	rm.SetArgs([]string{"source", "rm", "42"})
	if err := rm.Execute(); err != nil {
		t.Fatalf("rm: %v", err)
	}
}
