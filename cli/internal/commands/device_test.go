package commands

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDevice_List_Revoke(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/devices":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"id":"5","name":"macbook"}]`))
			return
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/devices/5":
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

	// list
	var out bytes.Buffer
	ls := NewRootCmd()
	ls.SetOut(&out)
	ls.SetArgs([]string{"device", "list"})
	if err := ls.Execute(); err != nil {
		t.Fatalf("device list: %v; out=%s", err, out.String())
	}
	if out.Len() == 0 {
		t.Fatalf("expected list output")
	}

	// revoke
	rv := NewRootCmd()
	rv.SetArgs([]string{"device", "revoke", "5"})
	if err := rv.Execute(); err != nil {
		t.Fatalf("device revoke: %v", err)
	}
}
