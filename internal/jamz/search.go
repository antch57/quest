package jamz

import "github.com/urfave/cli/v3"

func SearchCmd() *cli.Command {
	return &cli.Command{
		Name:     "search",
		Usage:    "search for shows in your area...",
		Commands: []*cli.Command{},
	}
}
