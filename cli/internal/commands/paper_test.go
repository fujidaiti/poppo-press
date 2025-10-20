package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaper_List_Read(t *testing.T) {
	// mock /v1/editions and /v1/editions/{id}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/editions":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]map[string]any{{"id": "17", "localDate": "2025-10-19", "articleCount": 2}})
			return
		case r.Method == http.MethodGet && r.URL.Path == "/v1/editions/17":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":        "17",
				"localDate": "2025-10-19",
				"articles": []map[string]any{
					{"position": 1, "article": map[string]any{"id": "202", "title": "A", "sourceId": "1"}},
					{"position": 2, "article": map[string]any{"id": "187", "title": "B", "sourceId": "2"}},
				},
			})
			return
		}
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	t.Cleanup(srv.Close)

	// init server and token via env login
	init := NewRootCmd()
	init.SetArgs([]string{"init", "--server", srv.URL})
	if err := init.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}
	t.Setenv("PP_TOKEN", "tok")
	lg := NewRootCmd()
	lg.SetArgs([]string{"login", "--device", "dev"})
	if err := lg.Execute(); err != nil {
		t.Fatalf("login: %v", err)
	}

	// list
	var out bytes.Buffer
	list := NewRootCmd()
	list.SetOut(&out)
	list.SetArgs([]string{"paper", "list", "--limit", "10", "--offset", "0"})
	if err := list.Execute(); err != nil {
		t.Fatalf("paper list: %v; out=%s", err, out.String())
	}
	if out.Len() == 0 {
		t.Fatalf("expected list output")
	}

	// read specific id 17
	out.Reset()
	read := NewRootCmd()
	read.SetOut(&out)
	read.SetArgs([]string{"paper", "read", "--id", "17"})
	if err := read.Execute(); err != nil {
		t.Fatalf("paper read: %v; out=%s", err, out.String())
	}
	if out.Len() == 0 {
		t.Fatalf("expected read output")
	}
}
