package jambase

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

func NewClient(apiKey string) *Client {
	return &Client{
		client:  &http.Client{Timeout: 10 * time.Second},
		apiKey:  apiKey,
		baseURL: "https://api.data.jambase.com/v3",
	}
}

// newRequest is a helper to handle the common headers for all calls
func (c *Client) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "quest/0.1")
	return req, nil
}
