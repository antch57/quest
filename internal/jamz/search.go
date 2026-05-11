package jamz

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/antch57/quest/internal/jamz/jambase"
	"github.com/antch57/quest/internal/tablefmt"
	"github.com/urfave/cli/v3"
)

// ErrApiKeyMissing is returned when JAMBASE_API_KEY is not set.
var ErrApiKeyMissing = errors.New("JAMBASE_API_KEY environment variable is required to search for shows")

// SearchOptions defines filters for looking up shows through the Jambase client.
type SearchOptions = jambase.SearchOptions

type showSearcher interface {
	SearchShows(ctx context.Context, opts SearchOptions) ([]jambase.Event, error)
}

const defaultJambaseBaseURL = "https://api.data.jambase.com/v3"

// SearchCmd returns the jamz search subcommand.
func SearchCmd() *cli.Command {
	return &cli.Command{
		Name:      "search",
		Usage:     "search for upcoming shows from Jambase...",
		UsageText: `quest jamz search --city "Denver"`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "city",
				Aliases:  []string{"c"},
				Usage:    "city to search around (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "country",
				Usage: "the two-letter ISO code for a country (optional, defaults to US)",
				Value: "US",
			},
			&cli.StringFlag{
				Name:    "artist",
				Aliases: []string{"a"},
				Usage:   "artist to filter by",
			},
			&cli.IntFlag{
				Name:    "radius",
				Aliases: []string{"r"},
				Usage:   "search radius in miles",
				Value:   25,
			},
			&cli.StringFlag{
				Name:    "date",
				Aliases: []string{"d"},
				Usage:   "date to search for shows (format: yyyy-mm-dd)",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"n"},
				Usage:   "maximum number of shows to return",
			},
			&cli.StringFlag{
				Name:    "venue",
				Aliases: []string{"v"},
				Usage:   "venue name to filter by",
			},
		},
		Action: runSearchCmd,
	}
}

func runSearchCmd(ctx context.Context, c *cli.Command) error {
	apiKey, err := apiKeyFromEnv()
	if err != nil {
		return err
	}

	client := jambaseClient(&http.Client{Timeout: 10 * time.Second}, apiKey, defaultJambaseBaseURL)
	opts := searchOptionsFromCommand(c)

	return searchAction(ctx, os.Stdout, client, opts)
}

func apiKeyFromEnv() (string, error) {
	apiKey := os.Getenv("JAMBASE_API_KEY")
	if apiKey == "" {
		return "", ErrApiKeyMissing
	}
	return apiKey, nil
}

func jambaseClient(httpClient *http.Client, apiKey string, baseURL string) *jambase.Client {
	return jambase.NewClient(
		httpClient,
		apiKey,
		baseURL,
	)
}

func searchOptionsFromCommand(c *cli.Command) SearchOptions {
	return SearchOptions{
		City:      c.String("city"),
		Country:   c.String("country"),
		Artist:    c.String("artist"),
		Radius:    c.Int("radius"),
		Limit:     c.Int("limit"),
		Date:      c.String("date"),
		VenueName: c.String("venue"),
	}
}

func searchAction(ctx context.Context, w io.Writer, s showSearcher, opts SearchOptions) error {
	result, err := s.SearchShows(ctx, opts)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		fmt.Fprintln(w, "no shows found")
		return nil
	}

	rows := make([][]string, 0, len(result))
	for _, event := range result {
		eventDate, startTime := formatEventDateAndStartTime(event.Date)
		rows = append(rows, []string{event.Name, eventDate, formatDoorTime(event.DoorTime), startTime, event.Venue, event.Address, event.Timezone})
	}

	tablefmt.Render(w, []string{"name", "date", "door time", "start time", "venue", "address", "timezone"}, rows)
	return nil
}

func formatEventDateAndStartTime(value string) (string, string) {
	if value == "" {
		return "-", "-"
	}

	layouts := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			if layout == "2006-01-02" {
				return parsed.Format("Mon Jan 2, 2006"), "-"
			}
			return parsed.Format("Mon Jan 2, 2006"), parsed.Format("3:04 PM")
		}
	}

	return value, "-"
}

func formatDoorTime(value string) string {
	if value == "" {
		return "-"
	}

	layouts := []string{"15:04:05", "15:04", time.RFC3339, "2006-01-02T15:04:05"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed.Format("3:04 PM")
		}
	}

	return value
}
