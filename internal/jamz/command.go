package jamz

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "jamz",
		Usage: "find some shows and save them to your quest logs...",
		Commands: []*cli.Command{
			SearchCmd(),
			SaveCmd(),
		},
	}
}
