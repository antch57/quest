package jamz

import "github.com/urfave/cli/v3"

func SaveCmd() *cli.Command {
	return &cli.Command{
		Name:     "save",
		Usage:    "save a show to your quest logs...",
		Commands: []*cli.Command{},
	}
}
