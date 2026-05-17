// Package jamz provides CLI commands for discovering and managing live shows.
package jamz

import (
	"github.com/urfave/cli/v3"
)

// Command returns the top-level jamz CLI command and its available subcommands.
func Command() *cli.Command {
	return &cli.Command{
		Name:  "jamz",
		Usage: "find some shows to go to...",
		Commands: []*cli.Command{
			SearchCmd(),
		},
	}
}
