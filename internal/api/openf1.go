package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

// Client is a minimal HTTP client for the OpenF1 API.
// It is safe for concurrent use.
type Client struct {
	baseURL *url.URL
	http    *http.Client
}

// NewClient creates a new OpenF1 client. If httpClient is nil, a default one is used.
func NewClient(base string, httpClient *http.Client) (*Client, error) {
	if base == "" {
		base = "https://openf1.org"
	}
	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{baseURL: u, http: httpClient}, nil
}

// GetCurrentSession fetches the sessions endpoint and returns the raw JSON response.
// The response is returned as bytes; callers are responsible for unmarshalling into their own types.
func (c *Client) GetCurrentSession(ctx context.Context) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, "sessions")
}

// GetRaceControlMessages fetches race control messages and returns raw JSON bytes.
func (c *Client) GetRaceControlMessages(ctx context.Context) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, path.Join("race-control", "messages"))
}

// doRequest executes an HTTP request against the configured base URL and returns the response body.
func (c *Client) doRequest(ctx context.Context, method, p string) ([]byte, error) {
	// Build URL relative to base
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, p)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("unexpected status: %d: %s", resp.StatusCode, string(body))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return b, nil
}
