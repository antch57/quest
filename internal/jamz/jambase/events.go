package jambase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ErrCityNotFound indicates that no matching city was found in Jambase geographies.
var ErrCityNotFound = errors.New("city not found")

// ErrInvalidLimit indicates that the requested limit was negative or exceeds the maximum allowed.
var ErrInvalidLimit = errors.New("limit must be between 0 and 50")

// ErrInvalidRadius indicates that the requested radius was negative.
var ErrInvalidRadius = errors.New("radius must be zero or greater")

// ErrInvalidDate indicates that a date did not match YYYY-MM-DD format.
var ErrInvalidDate = errors.New("date must use YYYY-MM-DD format")

// SearchOptions defines optional filters used when searching for events.
// All fields are optional; when empty, searches default to all upcoming shows in the US.
type SearchOptions struct {
	City      string
	Country   string
	Artist    string
	Date      string
	Limit     int
	Radius    int
	VenueName string
}

// SearchShows queries Jambase for events matching the provided search options.
// All search options are optional; when empty, returns all upcoming shows in the US.
// Results can be filtered by city, artist, venue, date, and radius (if city is provided).
func (c *Client) SearchShows(ctx context.Context, opts SearchOptions) ([]Event, error) {
	opts, err := validateSearchOptions(opts)
	if err != nil {
		return nil, err
	}

	var metroID string
	if opts.City != "" {
		// TODO: create local cache of metro IDs to avoid extra API call on every search by city
		metroID, err = c.cityToMetroID(ctx, opts.City, opts.Country)
		if err != nil {
			if errors.Is(err, ErrCityNotFound) {
				if opts.Country != "" {
					return nil, fmt.Errorf("city %q in country %q not found: %w", opts.City, opts.Country, err)
				}
				return nil, fmt.Errorf("city %q not found: %w", opts.City, err)
			}
			return nil, fmt.Errorf("lookup metro id: %w", err)
		}
	}

	reqURL, err := buildEventsURL(c.baseURL, metroID, opts)
	if err != nil {
		return nil, fmt.Errorf("build events request url: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodGet, reqURL)
	if err != nil {
		return nil, fmt.Errorf("build events request: %w", err)
	}

	resp, err := c.doRequestWithRetry(req)
	if err != nil {
		return nil, fmt.Errorf("send events request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search shows failed: status %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var response apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	events := make([]Event, 0, len(response.Events))
	for _, apiEvent := range response.Events {
		event := transformEvent(&apiEvent)
		events = append(events, event)
	}

	if opts.Limit > 0 && len(events) > opts.Limit {
		events = events[:opts.Limit]
	}

	return events, nil
}

// transformEvent converts an API event to a domain Event
func transformEvent(apiEvent *apiEvent) Event {
	return Event{
		Name:     apiEvent.Name,
		Date:     apiEvent.StartDate,
		DoorTime: apiEvent.DoorTime,
		Venue:    apiEvent.Location.Name,
		Address:  formatAddress(&apiEvent.Location.Address),
		Timezone: apiEvent.Location.Address.Timezone,
	}
}

// formatAddress constructs a human-readable address string from components
func formatAddress(addr *apiAddress) string {
	var parts []string

	if addr.StreetAddress != "" {
		parts = append(parts, addr.StreetAddress)
	}

	if addr.AddressLocality != "" {
		parts = append(parts, addr.AddressLocality)
	}

	if addr.AddressRegion.Name != "" {
		parts = append(parts, addr.AddressRegion.Name)
	}

	if addr.PostalCode != "" {
		parts = append(parts, addr.PostalCode)
	}

	if addr.AddressCountry.Name != "" {
		parts = append(parts, addr.AddressCountry.Name)
	}

	return strings.Join(parts, ", ")
}

func validateSearchOptions(opts SearchOptions) (SearchOptions, error) {
	opts.City = strings.TrimSpace(opts.City)
	opts.Country = strings.ToUpper(strings.TrimSpace(opts.Country))
	opts.Artist = strings.TrimSpace(opts.Artist)
	opts.Date = strings.TrimSpace(opts.Date)
	opts.VenueName = strings.TrimSpace(opts.VenueName)

	if opts.Limit < 0 || opts.Limit > 50 {
		return SearchOptions{}, fmt.Errorf("invalid limit %d: %w", opts.Limit, ErrInvalidLimit)
	}

	if opts.Radius < 0 {
		return SearchOptions{}, fmt.Errorf("invalid radius %d: %w", opts.Radius, ErrInvalidRadius)
	}

	if opts.Date != "" {
		if _, err := time.Parse("2006-01-02", opts.Date); err != nil {
			return SearchOptions{}, fmt.Errorf("invalid date %q: %w", opts.Date, ErrInvalidDate)
		}
	}

	return opts, nil
}

func buildEventsURL(baseURL, metroID string, opts SearchOptions) (string, error) {
	base, err := url.Parse(baseURL + "/events")
	if err != nil {
		return "", err
	}

	q := base.Query()
	if metroID != "" {
		q.Set("geoMetroId", metroID)
	}
	if opts.Artist != "" {
		q.Set("artistName", opts.Artist)
	}
	if opts.VenueName != "" {
		q.Set("venueName", opts.VenueName)
	}
	if opts.Radius > 0 {
		q.Set("geoRadiusAmount", strconv.Itoa(opts.Radius))
		q.Set("geoRadiusUnits", "mi")
	}
	if opts.Date != "" {
		q.Set("eventDateFrom", opts.Date)
		q.Set("eventDateTo", opts.Date)
	}
	if opts.Limit > 0 {
		q.Set("perPage", strconv.Itoa(opts.Limit))
	}
	base.RawQuery = q.Encode()

	return base.String(), nil
}

// cityToMetroID looks up the Jambase metro ID for a given city name
func (c *Client) cityToMetroID(ctx context.Context, city string, country string) (string, error) {
	reqURL, err := buildCityLookupURL(c.baseURL, city, country)
	if err != nil {
		return "", fmt.Errorf("build city lookup request url: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodGet, reqURL)
	if err != nil {
		return "", fmt.Errorf("build city lookup request: %w", err)
	}

	resp, err := c.doRequestWithRetry(req)
	if err != nil {
		return "", fmt.Errorf("send city lookup request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("city lookup failed: status %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var cityResp apiCity
	if err := json.NewDecoder(resp.Body).Decode(&cityResp); err != nil {
		return "", fmt.Errorf("failed to decode city lookup response: %w", err)
	}

	if !cityResp.Success || len(cityResp.Cities) == 0 {
		return "", ErrCityNotFound
	}

	// TODO: handle results with multiple matches (e.g. Denver, CO vs Denver, IA)
	return cityResp.Cities[0].Metro.Identifier, nil
}

func buildCityLookupURL(baseURL, city, country string) (string, error) {
	base, err := url.Parse(baseURL + "/geographies/cities")
	if err != nil {
		return "", err
	}

	q := base.Query()
	if city != "" {
		q.Set("geoCityName", city)
	}
	if country != "" {
		q.Set("geoCountryIso2", country)
	}
	base.RawQuery = q.Encode()

	return base.String(), nil
}
