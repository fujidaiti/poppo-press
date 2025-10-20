package commands

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLater_Add_List_Rm(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/read-later/202":
			w.WriteHeader(http.StatusNoContent)
			return
		case r.Method == http.MethodGet && r.URL.Path == "/v1/read-later":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"201"},{"id":"202"},{"id":"203"}]`))
			return
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/read-later/202":
			w.WriteHeader(http.StatusNoContent)
			return
		}
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	t.Cleanup(srv.Close)

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

	// add
	add := NewRootCmd()
	add.SetArgs([]string{"later", "add", "202"})
	if err := add.Execute(); err != nil {
		t.Fatalf("later add: %v", err)
	}

	// list
	var out bytes.Buffer
	ls := NewRootCmd()
	ls.SetOut(&out)
	ls.SetArgs([]string{"later", "list", "--limit", "1", "--offset", "1"})
	if err := ls.Execute(); err != nil {
		t.Fatalf("later list: %v; out=%s", err, out.String())
	}
	if out.Len() == 0 {
		t.Fatalf("expected list output")
	}

	// rm
	rm := NewRootCmd()
	rm.SetArgs([]string{"later", "rm", "202"})
	if err := rm.Execute(); err != nil {
		t.Fatalf("later rm: %v", err)
	}
}
