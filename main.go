// Command-line tool for organizing your day.
// quest is a lightweight, extensible CLI hub for todos, notes, and custom daily tasks.
// See https://github.com/antch57/quest for more information.
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
