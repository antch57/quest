package main

import (
	"context"
	"log"
	"os"

	"github.com/antch57/quest/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "quest",
		Usage: "get shit done..",
		Commands: []*cli.Command{
			cmd.CreateCmd(),
			cmd.ListCmd(),
			cmd.EditCmd(),
			cmd.DoneCmd(),
			cmd.DeleteCmd(),
			cmd.NukeCmd(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
