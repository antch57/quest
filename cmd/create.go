package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/antch57/quest/store"
	"github.com/urfave/cli/v3"
)

func CreateCmd() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create a new todo",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "due",
				Aliases: []string{"d"},
				Usage:   "due date of the todo (format: MM-DD-YYYY)",
			},
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Usage:   "project or folder for this todo",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			title := c.Args().First()
			if title == "" {
				return fmt.Errorf("usage: quest create <title>")
			}

			dueDateStr := c.String("due")
			var dueDate *time.Time
			if dueDateStr != "" {
				parsedDueDate, err := time.Parse("01-02-2006", dueDateStr)
				if err != nil {
					return fmt.Errorf("invalid due date format: %v", err)
				}
				dueDate = &parsedDueDate
			}

			project := c.String("project")

			todos, err := store.Load()
			if err != nil {
				return err
			}

			createID := func() int {
				max := 0
				for _, todo := range todos {
					if todo.ID > max {
						max = todo.ID
					}
				}
				return max + 1
			}

			todos = append(todos, store.Todo{
				ID:        createID(),
				Title:     title,
				CreatedAt: time.Now(),
				DueDate:   dueDate,
				Project:   project,
			})
			if err := store.Save(todos); err != nil {
				return err
			}
			fmt.Printf("Created todo: %s\n", title)
			return nil
		},
	}
}
