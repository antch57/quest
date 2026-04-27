package quest

import (
	"github.com/antch57/quest/internal/log"
	"github.com/urfave/cli/v3"
)

func App() *cli.Command {
	return &cli.Command{
		Name:  "quest",
		Usage: "get shit done..",
		Commands: []*cli.Command{
			log.Command(),
		},
	}
}
