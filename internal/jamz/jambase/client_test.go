// Package jambase provides a client for querying the Jambase API.
package jambase

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		client        *http.Client
		apiKey        string
		baseURL       string
		wantAPIKey    string
		wantBaseURL   string
		wantTimeout   time.Duration
		wantSamePtr   bool
		checkNonNilCl bool
	}{
		{
			name:        "uses provided client and base url",
			client:      &http.Client{Timeout: 3 * time.Second},
			apiKey:      "abc123",
			baseURL:     "https://example.test/api",
			wantAPIKey:  "abc123",
			wantBaseURL: "https://example.test/api",
			wantSamePtr: true,
		},
		{
			name:          "creates default client and base url when omitted",
			client:        nil,
			apiKey:        "key",
			baseURL:       "",
			wantAPIKey:    "key",
			wantBaseURL:   defaultBaseURL,
			wantTimeout:   10 * time.Second,
			checkNonNilCl: true,
		},
		{
			name:        "uses default base url when empty",
			client:      &http.Client{},
			apiKey:      "key-2",
			baseURL:     "",
			wantAPIKey:  "key-2",
			wantBaseURL: defaultBaseURL,
			wantSamePtr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.client, tt.apiKey, tt.baseURL)
			if got == nil {
				t.Fatalf("NewClient() returned nil")
			}

			if got.apiKey != tt.wantAPIKey {
				t.Errorf("NewClient().apiKey = %q, want %q", got.apiKey, tt.wantAPIKey)
			}

			if got.baseURL != tt.wantBaseURL {
				t.Errorf("NewClient().baseURL = %q, want %q", got.baseURL, tt.wantBaseURL)
			}

			if tt.wantSamePtr && got.client != tt.client {
				t.Errorf("NewClient().client did not preserve provided client pointer")
			}

			if tt.checkNonNilCl && got.client == nil {
				t.Fatalf("NewClient().client = nil, want non-nil")
			}

			if tt.wantTimeout > 0 && got.client.Timeout != tt.wantTimeout {
				t.Errorf("NewClient().client.Timeout = %v, want %v", got.client.Timeout, tt.wantTimeout)
			}

			if tt.wantTimeout == 0 && got.client == nil {
				t.Fatalf("NewClient().client = nil")
			}
		})
	}
}

func TestClient_newRequest(t *testing.T) {
	type fields struct {
		client  *http.Client
		apiKey  string
		baseURL string
	}
	type args struct {
		ctx    context.Context
		method string
		reqURL string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantAuth  string
		wantErr   bool
		errSubstr string
		ctxKey    any
		ctxValue  any
	}{
		{
			name: "valid request",
			fields: fields{
				client:  &http.Client{},
				apiKey:  "test-api-key",
				baseURL: "https://api.jambase.com",
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				reqURL: "https://api.jambase.com/events?artist=Phish",
			},
			wantAuth: "Bearer test-api-key",
			wantErr:  false,
		},
		{
			name: "includes context values",
			fields: fields{
				client:  &http.Client{},
				apiKey:  "ctx-key",
				baseURL: "https://api.jambase.com",
			},
			args: args{
				ctx:    context.WithValue(context.Background(), "trace-id", "abc-123"),
				method: http.MethodGet,
				reqURL: "https://api.jambase.com/events",
			},
			wantAuth: "Bearer ctx-key",
			ctxKey:   "trace-id",
			ctxValue: "abc-123",
		},
		{
			name: "invalid url returns wrapped error",
			fields: fields{
				client:  &http.Client{},
				apiKey:  "bad",
				baseURL: "https://api.jambase.com",
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				reqURL: "://bad-url",
			},
			wantErr:   true,
			errSubstr: "build request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				client:  tt.fields.client,
				apiKey:  tt.fields.apiKey,
				baseURL: tt.fields.baseURL,
			}
			got, err := c.newRequest(tt.args.ctx, tt.args.method, tt.args.reqURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Client.newRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if tt.errSubstr != "" && (err == nil || !strings.Contains(err.Error(), tt.errSubstr)) {
					t.Fatalf("Client.newRequest() error = %v, want substring %q", err, tt.errSubstr)
				}
				return
			}

			if got == nil {
				t.Fatalf("Client.newRequest() returned nil request")
			}

			if got.Method != tt.args.method {
				t.Errorf("Client.newRequest().Method = %q, want %q", got.Method, tt.args.method)
			}

			if got.URL.String() != tt.args.reqURL {
				t.Errorf("Client.newRequest().URL = %q, want %q", got.URL.String(), tt.args.reqURL)
			}

			if auth := got.Header.Get("Authorization"); auth != tt.wantAuth {
				t.Errorf("Client.newRequest().Authorization = %q, want %q", auth, tt.wantAuth)
			}

			if accept := got.Header.Get("Accept"); accept != "application/json" {
				t.Errorf("Client.newRequest().Accept = %q, want %q", accept, "application/json")
			}

			if ua := got.Header.Get("User-Agent"); ua != "quest/0.1" {
				t.Errorf("Client.newRequest().User-Agent = %q, want %q", ua, "quest/0.1")
			}

			if tt.ctxKey != nil {
				if v := got.Context().Value(tt.ctxKey); v != tt.ctxValue {
					t.Errorf("Client.newRequest() context value = %v, want %v", v, tt.ctxValue)
				}
			}
		})
	}
}

func TestClient_doRequestWithRetry(t *testing.T) {
	tests := []struct {
		name          string
		routes        map[string]testServerRoute
		retryAttempts int
		retryDelay    time.Duration
		wantStatus    int
		wantCalls     int32
		wantErr       bool
	}{
		{
			name: "returns success response without retry",
			routes: map[string]testServerRoute{
				"/events": {
					steps: []testServerStep{{statusCode: http.StatusOK, body: `{"ok":true}`}},
				},
			},
			retryAttempts: 3,
			retryDelay:    time.Millisecond,
			wantStatus:    http.StatusOK,
			wantCalls:     1,
			wantErr:       false,
		},
		{
			name: "retries and returns final error response status",
			routes: map[string]testServerRoute{
				"/events": {
					steps: []testServerStep{{statusCode: http.StatusBadGateway, body: `{"error":"upstream"}`}},
				},
			},
			retryAttempts: 3,
			retryDelay:    time.Millisecond,
			wantStatus:    http.StatusBadGateway,
			wantCalls:     3,
			wantErr:       false,
		},
		{
			name: "does not retry non-transient 4xx status",
			routes: map[string]testServerRoute{
				"/events": {
					steps: []testServerStep{{statusCode: http.StatusBadRequest, body: `{"error":"bad request"}`}},
				},
			},
			retryAttempts: 3,
			retryDelay:    time.Millisecond,
			wantStatus:    http.StatusBadRequest,
			wantCalls:     1,
			wantErr:       false,
		},
		{
			name: "retries after network error and then succeeds",
			routes: map[string]testServerRoute{
				"/events": {
					steps: []testServerStep{
						{disconnect: true},
						{statusCode: http.StatusOK, body: `{"ok":true}`},
					},
				},
			},
			retryAttempts: 3,
			retryDelay:    time.Millisecond,
			wantStatus:    http.StatusOK,
			wantCalls:     2,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, calls, cleanup := jambaseTestServer(t, tt.routes)
			defer cleanup()

			client.retryBaseDelay = tt.retryDelay
			client.retryMaxAttempts = tt.retryAttempts

			req, err := client.newRequest(context.Background(), http.MethodGet, client.baseURL+"/events")
			if err != nil {
				t.Fatalf("newRequest() error = %v", err)
			}

			resp, err := client.doRequestWithRetry(req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("doRequestWithRetry() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if resp == nil {
				t.Fatalf("doRequestWithRetry() returned nil response")
			}
			defer resp.Body.Close()

			if got := calls.Load(); got != tt.wantCalls {
				t.Fatalf("request count = %d, want %d", got, tt.wantCalls)
			}

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
