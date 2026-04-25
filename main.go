package main

import (
	"context"
	"log"
	"os"

	logcmd "github.com/antch57/quest/cmd/log"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "quest",
		Usage: "get shit done..",
		Commands: []*cli.Command{
			logcmd.Command(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
