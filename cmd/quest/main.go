// Package quest builds the top-level CLI command for the quest application.
package quest

import (
	"github.com/antch57/quest/internal/jamz"
	"github.com/antch57/quest/internal/log"
	"github.com/urfave/cli/v3"
)

// App returns the root CLI command for quest.
func App() *cli.Command {
	return &cli.Command{
		Name:  "quest",
		Usage: "get shit done..",
		Commands: []*cli.Command{
			log.Command(),
			jamz.Command(),
		},
	}
}
