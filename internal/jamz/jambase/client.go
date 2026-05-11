// Package jambase provides a client for querying the Jambase API.
package jambase

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.data.jambase.com/v3"

// Client wraps HTTP calls to the Jambase API.
type Client struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// NewClient creates a Jambase API client.
//
// If client is nil, a default client with a 10-second timeout is used.
// If baseURL is empty, the production Jambase v3 API URL is used.
func NewClient(client *http.Client, apiKey string, baseURL string) *Client {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Client{
		client:  client,
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// newRequest is a helper to handle the common headers for all calls
func (c *Client) newRequest(ctx context.Context, method, reqURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "quest/0.1")
	return req, nil
}
