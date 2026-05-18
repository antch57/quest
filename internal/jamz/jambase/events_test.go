package jambase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// jambaseTestServer starts a local httptest server with path-specific
// responses. The returned Client is already wired to it.
// Call the returned cleanup func (or defer it) to shut the server down.
type testServerStep struct {
	statusCode int
	body       string
	disconnect bool
}

type testServerRoute struct {
	statusCode int
	body       string
	steps      []testServerStep
	assert     func(t *testing.T, r *http.Request)
}

func jambaseTestServer(t *testing.T, routes map[string]testServerRoute) (*Client, *atomic.Int32, func()) {
	t.Helper()
	if routes == nil {
		routes = map[string]testServerRoute{}
	}

	pathCalls := &routeCallTracker{calls: make(map[string]int)}
	requestCount := &atomic.Int32{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		var step *testServerStep
		var statusCode int
		var body string

		w.Header().Set("Content-Type", "application/json")

		route, ok := routes[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"unhandled test route"}`))
			return
		}

		if route.assert != nil {
			route.assert(t, r)
		}

		if len(route.steps) > 0 {
			idx := pathCalls.next(r.URL.Path)
			if idx >= len(route.steps) {
				idx = len(route.steps) - 1
			}
			step = &route.steps[idx]
			statusCode = step.statusCode
			body = step.body
		} else {
			statusCode = route.statusCode
			body = route.body
		}

		if step != nil && step.disconnect {
			hj, ok := w.(http.Hijacker)
			if !ok {
				t.Fatalf("response writer does not support hijacking")
			}

			conn, _, err := hj.Hijack()
			if err != nil {
				t.Fatalf("hijack failed: %v", err)
			}
			_ = conn.Close()
			return
		}

		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}))
	c := &Client{
		client:  srv.Client(),
		apiKey:  "test-key",
		baseURL: srv.URL,
	}
	return c, requestCount, srv.Close
}

type routeCallTracker struct {
	mu    sync.Mutex
	calls map[string]int
}

func (r *routeCallTracker) next(path string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	current := r.calls[path]
	r.calls[path] = current + 1
	return current
}

// TODO: add tests for retry logic and caching behavior
func TestClient_SearchShows(t *testing.T) {
	type fields struct {
		apiKey string
	}
	type args struct {
		ctx  context.Context
		opts SearchOptions
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		routes      map[string]testServerRoute
		want        []Event
		wantErr     bool
		errContains string
		errIs       error
	}{
		{
			name:   "returns events on successful response",
			fields: fields{apiKey: "test-key"},
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					Artist: "goose",
					Limit:  1,
				},
			},
			routes: map[string]testServerRoute{
				"/events": {
					statusCode: http.StatusOK,
					body: `{
						"events": [
							{
								"@type": "Concert",
								"name": "Test Show 1",
								"startDate": "2024-07-04T20:00:00",
								"doorTime": "2024-07-04T19:00:00",
								"location": {
									"name": "Red Rocks Amphitheatre",
									"address": {
										"streetAddress": "18300 W Alameda Pkwy",
										"addressLocality": "Morrison",
										"addressRegion": {"name": "CO"},
										"postalCode": "80465",
										"x-timezone": "America/Denver",
										"addressCountry": {"name": "USA"}
									}
								}
							}
						]
					}`,
				},
			},
			want: []Event{
				{
					Name:     "Test Show 1",
					Date:     "2024-07-04T20:00:00",
					DoorTime: "2024-07-04T19:00:00",
					Venue:    "Red Rocks Amphitheatre",
					Address:  "18300 W Alameda Pkwy, Morrison, CO, 80465, USA",
					Timezone: "America/Denver",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid search options return error",
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					Limit: -1,
				},
			},
			wantErr: true,
		},
		{
			name: "filter by city returns events for that city",
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					City: "Denver",
				},
			},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body: `{
						"success": true,
						"cities": [
							{
								"name": "Denver",
								"identifier": "jambase:4227820",
								"containedInPlace": {
									"identifier": "jambase:8",
									"name": "Denver Area"
								}
							}
						]
					}`,
				},
				"/events": {
					statusCode: http.StatusOK,
					assert: func(t *testing.T, r *http.Request) {
						t.Helper()
						if got := r.URL.Query().Get("geoMetroId"); got != "jambase:8" {
							t.Fatalf("events geoMetroId = %q, want %q", got, "jambase:8")
						}
					},
					body: `{
						"events": [
							{
								"@type": "Concert",
								"name": "Denver Show",
								"startDate": "2024-07-04T20:00:00",
								"doorTime": "2024-07-04T19:00:00",
								"location": {
									"name": "Red Rocks Amphitheatre",
									"address": {
										"streetAddress": "18300 W Alameda Pkwy",
										"addressLocality": "Morrison",
										"addressRegion": {"name": "CO"},
										"postalCode": "80465",
										"x-timezone": "America/Denver",
										"addressCountry": {"name": "USA"}
									}
								}
							}
						]
					}`,
				},
			},
			want: []Event{
				{
					Name:     "Denver Show",
					Date:     "2024-07-04T20:00:00",
					DoorTime: "2024-07-04T19:00:00",
					Venue:    "Red Rocks Amphitheatre",
					Address:  "18300 W Alameda Pkwy, Morrison, CO, 80465, USA",
					Timezone: "America/Denver",
				},
			},
			wantErr: false,
		},
		{
			name: "city not found returns error",
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					City: "Nowhere",
				},
			},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body:       `{"success": false, "cities": []}`,
				},
			},
			wantErr:     true,
			errContains: fmt.Sprintf("city %q not found", "Nowhere"),
		},
		{
			name: "city not found with country returns error",
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					City:    "Nowhere",
					Country: "US",
				},
			},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body:       `{"success": false, "cities": []}`,
				},
			},
			wantErr:     true,
			errContains: fmt.Sprintf("city %q in country %q not found:", "Nowhere", "US"),
		},
		{
			name: "non-200 city lookup response returns error",
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					City: "Denver",
				},
			},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusInternalServerError,
					body:       ``,
				},
			},
			wantErr:     true,
			errContains: "lookup metro id:",
		},
		{
			name: "non-200 events response returns error",
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					City: "Denver",
				},
			},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body: `{
						"success": true,
						"cities": [
							{
								"name": "Denver",
								"identifier": "jambase:4227820",
								"containedInPlace": {
									"identifier": "jambase:8",
									"name": "Denver Area"
								}
							}
						]
					}`,
				},
				"/events": {
					statusCode: http.StatusInternalServerError,
					body:       ``,
				},
			},
			wantErr:     true,
			errContains: fmt.Sprintf("search shows failed: status %d", http.StatusInternalServerError),
		},
		{
			name: "malformed events response returns decode error",
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					Artist: "goose",
				},
			},
			routes: map[string]testServerRoute{
				"/events": {
					statusCode: http.StatusOK,
					body:       `{"events": [}`, // malformed JSON
				},
			},
			wantErr: true,
		},
		{
			name: "append events array on multiple pages of results",
			args: args{
				ctx: context.Background(),
				opts: SearchOptions{
					Artist: "goose",
					Limit:  1,
				},
			},
			routes: map[string]testServerRoute{
				"/events": {
					statusCode: http.StatusOK,
					body: `{
						"events": [
							{
								"name": "Goose Show",
								"startDate": "2024-07-04T20:00:00",
								"endDate": "2024-07-04T20:00:00",
								"doorTime": "2024-07-04T19:00:00",
								"location": {
									"name": "Red Rocks Amphitheatre",
									"address": {
										"streetAddress": "18300 W Alameda Pkwy",
										"addressLocality": "Morrison",
										"addressRegion": {"name": "CO"},
										"postalCode": "80465",
										"x-timezone": "America/Denver",
										"addressCountry": {"name": "USA"}
									}
								}
							},
							{
								"name": "Goose Show 2",
								"startDate": "2024-07-05T20:00:00",
								"endDate": "2024-07-05T20:00:00",
								"doorTime": "2024-07-05T19:00:00",
								"location": {
									"name": "Red Rocks Amphitheatre",
									"address": {
										"streetAddress": "18300 W Alameda Pkwy",
										"addressLocality": "Morrison",
										"addressRegion": {"name": "CO"},
										"postalCode": "80465",
										"x-timezone": "America/Denver",
										"addressCountry": {"name": "USA"}
									}
								}
							}
						]
					}`,
				},
			},
			want: []Event{
				{
					Name:     "Goose Show",
					Date:     "2024-07-04T20:00:00",
					DoorTime: "2024-07-04T19:00:00",
					Venue:    "Red Rocks Amphitheatre",
					Address:  "18300 W Alameda Pkwy, Morrison, CO, 80465, USA",
					Timezone: "America/Denver",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _, cleanup := jambaseTestServer(t, tt.routes)
			defer cleanup()
			c.apiKey = tt.fields.apiKey

			got, err := c.SearchShows(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Client.SearchShows() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if tt.errContains != "" && (err == nil || !strings.Contains(err.Error(), tt.errContains)) {
					t.Fatalf("SearchShows() error = %q, want to contain %q", err, tt.errContains)
				}
				if tt.errIs != nil && !errors.Is(err, tt.errIs) {
					t.Fatalf("SearchShows() error = %v, want errors.Is(..., %v)", err, tt.errIs)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SearchShows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformEvent(t *testing.T) {
	type args struct {
		apiEvent *apiEvent
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			name: "transforms apiEvent to Event correctly",
			args: args{
				apiEvent: &apiEvent{
					Type:      "Concert",
					Name:      "Test Show",
					StartDate: "2024-07-04T20:00:00",
					EndDate:   "2024-07-04T20:00:00",
					DoorTime:  "2024-07-04T19:00:00",
					Location: apiVenue{
						Name: "Red Rocks Amphitheatre",
						Address: apiAddress{
							StreetAddress:   "18300 W Alameda Pkwy",
							AddressLocality: "Morrison",
							AddressRegion: apiNamedLocation{
								Name: "CO",
							},
							PostalCode: "80465",
							Timezone:   "America/Denver",
							AddressCountry: apiNamedLocation{
								Name: "USA",
							},
						},
					},
				},
			},
			want: Event{
				Name:     "Test Show",
				Date:     "2024-07-04T20:00:00",
				DoorTime: "2024-07-04T19:00:00",
				Venue:    "Red Rocks Amphitheatre",
				Address:  "18300 W Alameda Pkwy, Morrison, CO, 80465, USA",
				Timezone: "America/Denver",
			},
		},
		{
			name: "handles missing optional fields",
			args: args{
				apiEvent: &apiEvent{
					Type:      "Concert",
					Name:      "Test Show",
					StartDate: "2024-07-04T20:00:00",
					Location: apiVenue{
						Name: "Red Rocks Amphitheatre",
						Address: apiAddress{
							Timezone: "America/Denver",
						},
					},
				},
			},
			want: Event{
				Name:     "Test Show",
				Date:     "2024-07-04T20:00:00",
				DoorTime: "",
				Venue:    "Red Rocks Amphitheatre",
				Address:  "",
				Timezone: "America/Denver",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformEvent(tt.args.apiEvent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatAddress(t *testing.T) {
	type args struct {
		addr *apiAddress
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "formats full address correctly",
			args: args{
				addr: &apiAddress{
					StreetAddress:   "123 Main St",
					AddressLocality: "Denver",
					PostalCode:      "80202",
					Timezone:        "America/Denver",
					AddressRegion: apiNamedLocation{
						Name: "CO",
					},
					AddressCountry: apiNamedLocation{
						Name: "USA",
					},
				},
			},
			want: "123 Main St, Denver, CO, 80202, USA",
		},
		{
			name: "formats partial address correctly",
			args: args{
				addr: &apiAddress{
					AddressLocality: "Denver",
					Timezone:        "America/Denver",
				},
			},
			want: "Denver",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatAddress(tt.args.addr); got != tt.want {
				t.Errorf("formatAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateSearchOptions(t *testing.T) {
	type args struct {
		opts SearchOptions
	}
	tests := []struct {
		name    string
		args    args
		want    SearchOptions
		wantErr bool
	}{
		{
			name: "valid options",
			args: args{
				opts: SearchOptions{
					City:    "New York",
					Country: "US",
					Radius:  50,
				},
			},
			want: SearchOptions{
				City:    "New York",
				Country: "US",
				Radius:  50,
			},
			wantErr: false,
		},
		{
			name: "invalid limit too low",
			args: args{
				opts: SearchOptions{
					Country: "US",
					Limit:   -5,
				},
			},
			want:    SearchOptions{},
			wantErr: true,
		},
		{
			name: "invalid limit too high",
			args: args{
				opts: SearchOptions{
					Country: "US",
					Limit:   150,
				},
			},
			want:    SearchOptions{},
			wantErr: true,
		},
		{
			name: "invalid radius",
			args: args{
				opts: SearchOptions{
					Country: "US",
					Radius:  -10,
				},
			},
			want:    SearchOptions{},
			wantErr: true,
		},
		{
			name: "invalid date",
			args: args{
				opts: SearchOptions{
					Country: "US",
					Date:    "2024-13-01",
				},
			},
			want:    SearchOptions{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateSearchOptions(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateSearchOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateSearchOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildEventsURL(t *testing.T) {
	type args struct {
		baseURL string
		metroID string
		opts    SearchOptions
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "build URL with metroID and artist",
			args: args{
				baseURL: "https://api.jambase.com",
				metroID: "jambase:123",
				opts: SearchOptions{
					Artist: "goose",
				},
			},
			want:    "https://api.jambase.com/events?artistName=goose&geoMetroId=jambase%3A123",
			wantErr: false,
		},
		{
			name: "build URL with no search options or metroID",
			args: args{
				baseURL: "https://api.jambase.com",
				metroID: "",
				opts:    SearchOptions{},
			},
			want:    "https://api.jambase.com/events",
			wantErr: false,
		},
		{
			name: "invalid base URL returns error",
			args: args{
				baseURL: "://bad-url",
				metroID: "",
				opts:    SearchOptions{},
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "build URL with all search options",
			args: args{
				baseURL: "https://api.jambase.com",
				metroID: "jambase:123",
				opts: SearchOptions{
					Artist:    "goose",
					VenueName: "red rocks",
					Date:      "2024-07-04",
					Radius:    25,
					Limit:     10,
				},
			},
			want:    "https://api.jambase.com/events?artistName=goose&eventDateFrom=2024-07-04&eventDateTo=2024-07-04&geoMetroId=jambase%3A123&geoRadiusAmount=25&geoRadiusUnits=mi&perPage=10&venueName=red+rocks",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildEventsURL(tt.args.baseURL, tt.args.metroID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("buildEventsURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("buildEventsURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_cityToMetroID(t *testing.T) {
	type fields struct {
		apiKey string
	}
	type args struct {
		ctx     context.Context
		city    string
		country string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		routes  map[string]testServerRoute
		want    string
		wantErr bool
		baseURL string // overrides the test server URL when set
	}{
		{
			name:   "returns metro id for matched city",
			fields: fields{apiKey: "test-key"},
			args:   args{ctx: context.Background(), city: "Denver", country: "US"},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body: `{
						"success": true,
						"cities": [
							{
								"name": "Denver",
								"identifier": "jambase:4227820",
								"containedInPlace": {
									"identifier": "jambase:8",
									"name": "Denver Area"
								}
							}
						]
					}`,
				},
			},
			want: "jambase:8",
		},
		{
			name:   "returns ErrCityNotFound when success false",
			fields: fields{apiKey: "test-key"},
			args:   args{ctx: context.Background(), city: "Nowhere", country: ""},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body:       `{"success": false, "cities": []}`,
				},
			},
			wantErr: true,
		},
		{
			name:   "returns ErrCityNotFound when cities empty",
			fields: fields{apiKey: "test-key"},
			args:   args{ctx: context.Background(), city: "Ghost Town", country: ""},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body:       `{"success": true, "cities": []}`,
				},
			},
			wantErr: true,
		},
		{
			name:   "wraps error on non-200 response",
			fields: fields{apiKey: "test-key"},
			args:   args{ctx: context.Background(), city: "Denver", country: ""},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusInternalServerError,
					body:       ``,
				},
			},
			wantErr: true,
		},
		{
			name:   "bad url returns error",
			fields: fields{apiKey: "test-key"},
			args: args{
				ctx:     context.Background(),
				city:    "Denver",
				country: "",
			},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body:       `{"success": true, "cities": []}`,
				},
			},
			wantErr: true,
			baseURL: "://invalid",
		},
		{
			name:   "malformed json response returns error",
			fields: fields{apiKey: "test-key"},
			args: args{
				ctx:     context.Background(),
				city:    "Denver",
				country: "",
			},
			routes: map[string]testServerRoute{
				"/geographies/cities": {
					statusCode: http.StatusOK,
					body:       `{"success": true, "cities": [}`,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _, cleanup := jambaseTestServer(t, tt.routes)
			defer cleanup()
			c.apiKey = tt.fields.apiKey
			if tt.baseURL != "" {
				c.baseURL = tt.baseURL
			}

			got, err := c.cityToMetroID(tt.args.ctx, tt.args.city, tt.args.country)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Client.cityToMetroID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("Client.cityToMetroID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildCityLookupURL(t *testing.T) {
	type args struct {
		baseURL string
		city    string
		country string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "build URL with city and country",
			args: args{
				baseURL: "https://api.jambase.com",
				city:    "Denver",
				country: "US",
			},
			want:    "https://api.jambase.com/geographies/cities?geoCityName=Denver&geoCountryIso2=US",
			wantErr: false,
		},
		{
			name: "build URL with no city",
			args: args{
				baseURL: "https://api.jambase.com",
				city:    "",
				country: "US",
			},
			want:    "https://api.jambase.com/geographies/cities?geoCountryIso2=US",
			wantErr: false,
		},
		{
			name: "build URL with no country",
			args: args{
				baseURL: "https://api.jambase.com",
				city:    "Denver",
				country: "",
			},
			want:    "https://api.jambase.com/geographies/cities?geoCityName=Denver",
			wantErr: false,
		},
		{
			name: "invalid base URL returns error",
			args: args{
				baseURL: "://bad-url",
				city:    "Denver",
				country: "US",
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildCityLookupURL(tt.args.baseURL, tt.args.city, tt.args.country)
			if (err != nil) != tt.wantErr {
				t.Fatalf("buildCityLookupURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("buildCityLookupURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
