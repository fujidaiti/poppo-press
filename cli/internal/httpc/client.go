package httpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Option func(*Client)

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option { return func(c *Client) { c.http.Timeout = d } }

// WithVerbose enables verbose tracing to the provided writer with redacted secrets.
func WithVerbose(w io.Writer) Option { return func(c *Client) { c.verbose = w } }

// Error represents an operation error with a mapped exit code per CLI policy.
type Error struct {
	Code int   // 0 ok; 1 generic; 2 validation; 3 auth; 4 network
	Err  error // underlying error
}

func (e *Error) Error() string { return e.Err.Error() }
func (e *Error) Unwrap() error { return e.Err }

type Client struct {
	base    *url.URL
	token   string
	http    *http.Client
	verbose io.Writer
}

// New constructs a client with base URL and optional bearer token.
func New(baseURL string, token string, opts ...Option) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{base: u, token: token, http: &http.Client{Timeout: 10 * time.Second}}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

// NewRequest builds an HTTP request joined to the base URL and injects headers.
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := c.base.ResolveReference(rel)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return req, nil
}

// Do performs the HTTP request, mapping errors and optionally printing verbose traces.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.verbose != nil {
		fmt.Fprintf(c.verbose, ">>> %s %s\n", req.Method, req.URL.String())
		// redact Authorization header
		if v := req.Header.Get("Authorization"); v != "" {
			fmt.Fprintf(c.verbose, "Authorization: Bearer ****\n")
		}
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &Error{Code: 4, Err: err}
	}
	if c.verbose != nil {
		fmt.Fprintf(c.verbose, "<<< %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	}
	// Drain body for completeness; caller shouldn't rely on it on errors
	var buf bytes.Buffer
	_, _ = io.CopyN(&buf, resp.Body, 1024)
	resp.Body.Close()
	code := mapStatusToExitCode(resp.StatusCode)
	msg := strings.TrimSpace(buf.String())
	if msg == "" {
		msg = resp.Status
	}
	return nil, &Error{Code: code, Err: errors.New(msg)}
}

func mapStatusToExitCode(status int) int {
	switch status {
	case 400, 422:
		return 2
	case 401, 403:
		return 3
	default:
		if status >= 500 {
			return 1
		}
		return 1
	}
}
