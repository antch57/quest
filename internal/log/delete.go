package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/antch57/quest/internal/store"
	"github.com/urfave/cli/v3"
)

func deleteAction(id string) error {
	todos, idx, err := store.LoadAndFindIndexByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("task with id %s not found", id)
		}
		return err
	}
	todos[idx].Deleted = true
	fmt.Printf("you have deleted: \"%s\"\n", todos[idx].Title)
	if err := store.Save(todos); err != nil {
		return err
	}
	return nil
}

func DeleteCmd() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete task by id",
		UsageText: "quest log delete --id <task id>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "id of the task to delete (required)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			id := c.String("id")
			err := deleteAction(id)
			if err != nil {
				cli.ShowCommandHelp(ctx, c, "delete")
			}
			return err
		},
	}
}

func nukeAction(r io.Reader) error {
	fmt.Print("are you sure you want to nuke all tasks? (y/n): ")
	var response string
	fmt.Fscan(r, &response)

	if strings.ToLower(response) != "y" {
		fmt.Println("aborted.")
		return nil
	}

	if err := store.Nuke(); err != nil {
		return err
	}

	fmt.Println("you have nuked all tasks")
	return nil
}

func NukeCmd() *cli.Command {
	return &cli.Command{
		Name:      "nuke",
		Usage:     "delete .quest/todo.json file",
		UsageText: "quest log nuke",
		Action: func(ctx context.Context, c *cli.Command) error {
			err := nukeAction(os.Stdin)
			if err != nil {
				cli.ShowCommandHelp(ctx, c, "nuke")
			}
			return err
		},
	}
}
