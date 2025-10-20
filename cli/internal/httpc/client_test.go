package httpc

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBaseURLAndTokenInjection(t *testing.T) {
	var gotAuth string
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(srv.Close)

	c, err := New(srv.URL, "abc", WithTimeout(5*time.Second))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	req, err := c.NewRequest(context.Background(), http.MethodGet, "/v1/x", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	_, err = c.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if gotAuth != "Bearer abc" {
		t.Fatalf("auth header = %q, want %q", gotAuth, "Bearer abc")
	}
	if gotPath != "/v1/x" {
		t.Fatalf("path = %q, want %q", gotPath, "/v1/x")
	}
}

func TestVerboseRedaction(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(srv.Close)

	var buf bytes.Buffer
	c, err := New(srv.URL, "sensitive-token", WithVerbose(&buf), WithTimeout(2*time.Second))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	req, err := c.NewRequest(context.Background(), http.MethodGet, "/x", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	if _, err := c.Do(req); err != nil {
		t.Fatalf("Do: %v", err)
	}
	s := buf.String()
	if bytes.Contains([]byte(s), []byte("sensitive-token")) {
		t.Fatalf("trace contains raw token: %s", s)
	}
	if !bytes.Contains([]byte(s), []byte("Authorization: Bearer ****")) {
		t.Fatalf("trace missing redacted Authorization header: %s", s)
	}
}

func TestExitCodeMapping(t *testing.T) {
	cases := []struct {
		status   int
		wantCode int
	}{
		{400, 2},
		{422, 2},
		{401, 3},
		{403, 3},
		{500, 1},
	}
	for _, tc := range cases {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
		}))
		c, err := New(srv.URL, "")
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		req, _ := c.NewRequest(context.Background(), http.MethodGet, "/x", nil)
		_, err = c.Do(req)
		if err == nil {
			t.Fatalf("expected error for status %d", tc.status)
		}
		var he *Error
		if !errors.As(err, &he) {
			t.Fatalf("expected *Error, got %T: %v", err, err)
		}
		if he.Code != tc.wantCode {
			t.Fatalf("status %d -> code %d, want %d", tc.status, he.Code, tc.wantCode)
		}
		srv.Close()
	}
}

func TestExitCodeMapping_Network(t *testing.T) {
	// Unreachable port likely yields connection refused
	c, err := New("http://127.0.0.1:1", "", WithTimeout(500*time.Millisecond))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	req, _ := c.NewRequest(context.Background(), http.MethodGet, "/x", nil)
	_, err = c.Do(req)
	if err == nil {
		t.Fatalf("expected network error")
	}
	var he *Error
	if !errors.As(err, &he) {
		t.Fatalf("expected *Error, got %T: %v", err, err)
	}
	if he.Code != 4 {
		t.Fatalf("network error -> code %d, want 4", he.Code)
	}
}
