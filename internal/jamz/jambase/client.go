// Package jambase provides a client for querying the Jambase API.
package jambase

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.data.jambase.com/v3"

const (
	defaultRetryMaxAttempts = 3
	defaultRetryBaseDelay   = 300 * time.Millisecond
)

// Client wraps HTTP calls to the Jambase API.
type Client struct {
	client  *http.Client
	apiKey  string
	baseURL string

	retryMaxAttempts int
	retryBaseDelay   time.Duration
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

		retryMaxAttempts: defaultRetryMaxAttempts,
		retryBaseDelay:   defaultRetryBaseDelay,
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

func (c *Client) doRequestWithRetry(req *http.Request) (*http.Response, error) {
	attempts := c.retryMaxAttempts
	if attempts <= 0 {
		attempts = defaultRetryMaxAttempts
	}

	baseDelay := c.retryBaseDelay
	if baseDelay <= 0 {
		baseDelay = defaultRetryBaseDelay
	}

	for attempt := 1; attempt <= attempts; attempt++ {
		resp, err := c.client.Do(req)
		if err == nil && resp.StatusCode < http.StatusBadRequest {
			return resp, nil
		}

		shouldRetry := err != nil
		if err == nil {
			shouldRetry = isRetryableStatusCode(resp.StatusCode)
		}

		if !shouldRetry || attempt == attempts {
			if err != nil {
				if resp != nil {
					resp.Body.Close()
				}
				return nil, err
			}
			return resp, nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		if err := sleepWithContext(req.Context(), baseDelay); err != nil {
			return nil, err
		}

		baseDelay *= 2
	}

	return nil, fmt.Errorf("request failed after %d attempts", attempts)
}

func isRetryableStatusCode(statusCode int) bool {
	if statusCode >= http.StatusInternalServerError {
		return true
	}

	switch statusCode {
	case http.StatusRequestTimeout, http.StatusTooEarly, http.StatusTooManyRequests:
		return true
	default:
		return false
	}
}

func sleepWithContext(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
