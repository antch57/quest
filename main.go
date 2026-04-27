package main

import (
	"context"
	"log"
	"os"

	"github.com/antch57/quest/cmd/quest"
)

func main() {
	app := quest.App()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
