package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/antch57/quest/internal/store"
	"github.com/urfave/cli/v3"
)

func DoneCmd() *cli.Command {
	return &cli.Command{
		Name:      "done",
		Usage:     "mark task as done by id",
		UsageText: "quest log done --id <task id>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "id of the task to mark as done",
				Required: true,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			id := c.String("id")
			todos, idx, err := store.LoadAndFindIndexByID(id)
			if err != nil {
				if err == os.ErrNotExist {
					return fmt.Errorf("task with id %s not found", id)
				}
				return err
			}
			todos[idx].Done = true
			fmt.Printf("you have completed: \"%s\"\n", todos[idx].Title)
			return store.Save(todos)
		},
	}
}

func EditCmd() *cli.Command {
	return &cli.Command{
		Name:      "edit",
		Usage:     "edit a task by id",
		UsageText: "quest log edit --id <task id> [options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "id of the task to edit (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "new title for the todo",
			},
			&cli.StringFlag{
				Name:    "due",
				Aliases: []string{"d"},
				Usage:   "new due date for the todo (format: mm-dd-yyyy)",
			},
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Usage:   "new project for the todo",
			},
			&cli.BoolFlag{
				Name:  "clear-due",
				Usage: "clear the due date for the todo (no value needed)",
			},
			&cli.BoolFlag{
				Name:  "done",
				Usage: "mark the todo as done",
			},
			&cli.BoolFlag{
				Name:  "undone",
				Usage: "mark the todo as not done",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			id := c.String("id")
			todos, idx, err := store.LoadAndFindIndexByID(id)
			if err != nil {
				if err == os.ErrNotExist {
					return fmt.Errorf("task with id %s not found", id)
				}
				return err
			}

			if c.Bool("done") && c.Bool("undone") {
				return fmt.Errorf("cannot use both --done and --undone flags at the same time")
			}

			if c.Bool("clear-due") && c.IsSet("due") {
				return fmt.Errorf("cannot use both --clear-due and --due flags at the same time")
			}

			if c.IsSet("title") {
				todos[idx].Title = c.String("title")
			}
			if c.Bool("clear-due") {
				todos[idx].DueDate = nil
			}
			if c.IsSet("due") {
				dueDateStr := c.String("due")
				parsedDueDate, err := time.Parse("01-02-2006", dueDateStr)
				if err != nil {
					return fmt.Errorf("invalid due date format: %v", err)
				}
				todos[idx].DueDate = &parsedDueDate
			}
			if c.Bool("done") {
				todos[idx].Done = true
			}
			if c.Bool("undone") {
				todos[idx].Done = false
			}
			if c.IsSet("project") {
				todos[idx].Project = c.String("project")
			}

			if err := store.Save(todos); err != nil {
				return err
			}

			fmt.Printf("you have updated: \"%s\"\n", todos[idx].Title)

			return nil
		},
	}
}
