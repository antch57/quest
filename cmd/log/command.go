package log

import (
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "log",
		Usage: "manage your quest log...",
		Commands: []*cli.Command{
			CreateCmd(),
			ListCmd(),
			EditCmd(),
			DoneCmd(),
			DeleteCmd(),
			NukeCmd(),
		},
	}
}
