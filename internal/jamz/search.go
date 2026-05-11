package jamz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/antch57/quest/internal/jamz/jambase"
	"github.com/urfave/cli/v3"
)

var ErrApiKeyMissing = errors.New("JAMBASE_API_KEY environment variable is required to search for shows")

type SearchOptions = jambase.SearchOptions

func searchAction(ctx context.Context, opts SearchOptions) error {
	apiKey := os.Getenv("JAMBASE_API_KEY")
	if apiKey == "" {
		return ErrApiKeyMissing
	}

	client := jambase.NewClient(apiKey)
	result, err := jambase.SearchShows(ctx, client, opts)
	if err != nil {
		return err
	}

	formatted, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(formatted))
	return nil
}

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
				Value:   10,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			err := searchAction(ctx, SearchOptions{
				City:   c.String("city"),
				Artist: c.String("artist"),
				Radius: c.Int("radius"),
				Limit:  c.Int("limit"),
				Date:   c.String("date"),
			})
			if err != nil {
				if errors.Is(err, ErrApiKeyMissing) {
					return fmt.Errorf("api key not found: %w", err)
				}
				cli.ShowCommandHelp(ctx, c, "search")
			}

			return err
		},
	}
}
