package jambase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var ErrCityNotFound = errors.New("city not found")

type SearchOptions struct {
	City    string
	Country string
	Artist  string
	Date    string
	Limit   int
	Radius  int
}

func SearchShows(ctx context.Context, c *Client, opts SearchOptions) ([]Event, error) {
	// TODO: this is a bit of a hack - the API requires a metro ID for searching, but users will want to search by city name. We should probably cache metro IDs locally to avoid hitting the API every time, but for now we'll just do a lookup on each search.
	metroID, err := cityToMetroID(ctx, c, opts.City, "")
	if err != nil {
		if errors.Is(err, ErrCityNotFound) {
			return nil, fmt.Errorf("city '%s' not found: %w", opts.City, err)
		}
		return nil, err
	}

	base, err := url.Parse(c.baseURL + "/events")
	if err != nil {
		return nil, err
	}

	q := base.Query()
	q.Set("geoMetroId", metroID)
	if opts.Artist != "" {
		q.Set("artistName", opts.Artist)
	}
	base.RawQuery = q.Encode()

	req, err := c.newRequest(ctx, http.MethodGet, base.String())
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: here for testing so i don't have to hit the API every time, but should be replaced with a real request
	// fixture, err := os.Open("internal/jamz/jambase/sample-concert-event.json")
	// fixture, err := os.Open("internal/jamz/jambase/testData/jambase-event-response.json")
	// if err != nil {
	// 	return nil, err
	// }
	// defer fixture.Close()

	var response apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	events := make([]Event, 0, len(response.Events))
	for _, apiEvent := range response.Events {
		event := transformEvent(&apiEvent)
		events = append(events, event)
	}

	return events, nil
}

// transformEvent converts an API event to a domain Event
func transformEvent(apiEvent *apiEvent) Event {
	return Event{
		Name:     apiEvent.Name,
		DoorTime: apiEvent.DoorTime,
		Venue:    apiEvent.Location.Name,
		Address:  formatAddress(&apiEvent.Location.Address),
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

// cityToMetroID looks up the Jambase metro ID for a given city name
func cityToMetroID(ctx context.Context, c *Client, city string, country string) (string, error) {
	base, err := url.Parse(c.baseURL + "/geographies/cities")
	if err != nil {
		return "", err
	}

	q := base.Query()
	q.Set("geoCityName", city)
	if country != "" {
		q.Set("geoCountryIso2", country)
	}
	base.RawQuery = q.Encode()

	req, err := c.newRequest(ctx, http.MethodGet, base.String())
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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
