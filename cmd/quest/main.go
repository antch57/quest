// Package quest builds the top-level CLI command for the quest application.
package quest

import (
	"runtime/debug"

	"github.com/antch57/quest/internal/jamz"
	"github.com/antch57/quest/internal/log"
	"github.com/urfave/cli/v3"
)

var version = "dev"

// App returns the root CLI command for quest.
func App() *cli.Command {
	return &cli.Command{
		Name:    "quest",
		Usage:   "a cli tool to help you manage your quests",
		Version: resolvedVersion(),
		Commands: []*cli.Command{
			log.Command(),
			jamz.Command(),
		},
	}
}

func resolvedVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" {
		return "devel"
	}
	return info.Main.Version
}
