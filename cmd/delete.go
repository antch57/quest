package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/antch57/quest/store"
	"github.com/urfave/cli/v3"
)

func DeleteCmd() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "delete task by id",
		Action: func(ctx context.Context, c *cli.Command) error {
			arg := c.Args().First()
			if arg == "" {
				return fmt.Errorf("usage: quest delete <task id>")
			}

			id, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}

			todos, idx, err := store.LoadAndFindIndexByID(id)
			if err != nil {
				if err == os.ErrNotExist {
					return fmt.Errorf("task with id %d not found", id)
				}
				return err
			}

			todos[idx].Deleted = true
			fmt.Printf("you have deleted: \"%s\"\n", todos[idx].Title)

			return store.Save(todos)
		},
	}
}

func NukeCmd() *cli.Command {
	return &cli.Command{
		Name:  "nuke",
		Usage: "delete .quest/todo.json file",
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Print("are you sure you want to nuke all tasks? (y/n): ")
			var response string
			fmt.Scanln(&response)

			if strings.ToLower(response) != "y" {
				fmt.Println("aborted.")
				return nil
			}

			if err := store.Nuke(); err != nil {
				return err
			}

			fmt.Println("you have nuked all tasks")
			return nil
		},
	}
}
