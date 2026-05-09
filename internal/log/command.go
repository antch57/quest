// Package log contains quest todo command handlers.
//
// It defines the "quest log" command tree and the behavior for creating,
// listing, editing, completing, deleting, and nuking todos.
package log

import (
	"github.com/urfave/cli/v3"
)

// Command returns the top-level "log" command and all nested subcommands.
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
