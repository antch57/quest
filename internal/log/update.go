package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/antch57/quest/internal/store"
	"github.com/urfave/cli/v3"
)

type EditOptions struct {
	ID       string
	Title    *string
	Due      *string
	Project  *string
	ClearDue bool
	Done     bool
	Undone   bool
}

func doneAction(w io.Writer, id string) error {
	todos, idx, err := store.LoadAndFindIndexByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("task with id %s not found", id)
		}
		return err
	}
	todos[idx].Done = true
	fmt.Fprintf(w, "you have completed: \"%s\"\n", todos[idx].Title)
	if err := store.Save(todos); err != nil {
		return err
	}
	return nil
}

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
			err := doneAction(os.Stdout, id)
			if err != nil {
				cli.ShowCommandHelp(ctx, c, "done")
			}
			return err
		},
	}
}

func editAction(w io.Writer, opts EditOptions) error {
	todos, idx, err := store.LoadAndFindIndexByID(opts.ID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("task with id %s not found", opts.ID)
		}
		return err
	}

	if opts.Done && opts.Undone {
		return fmt.Errorf("cannot use both --done and --undone flags at the same time")
	}
	if opts.ClearDue && opts.Due != nil {
		return fmt.Errorf("cannot use both --clear-due and --due flags at the same time")
	}

	if opts.Title != nil {
		todos[idx].Title = *opts.Title
	}
	if opts.ClearDue {
		todos[idx].DueDate = nil
	}
	if opts.Due != nil {
		parsedDueDate, err := time.Parse("01-02-2006", *opts.Due)
		if err != nil {
			return fmt.Errorf("invalid due date format: %v", err)
		}
		todos[idx].DueDate = &parsedDueDate
	}
	if opts.Done {
		todos[idx].Done = true
	}
	if opts.Undone {
		todos[idx].Done = false
	}
	if opts.Project != nil {
		todos[idx].Project = *opts.Project
	}

	if err := store.Save(todos); err != nil {
		return err
	}

	fmt.Fprintf(w, "you have updated: \"%s\"\n", todos[idx].Title)
	return nil
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
			opts := EditOptions{
				ID:       c.String("id"),
				ClearDue: c.Bool("clear-due"),
				Done:     c.Bool("done"),
				Undone:   c.Bool("undone"),
			}
			if c.IsSet("title") {
				title := c.String("title")
				opts.Title = &title
			}
			if c.IsSet("due") {
				due := c.String("due")
				opts.Due = &due
			}
			if c.IsSet("project") {
				project := c.String("project")
				opts.Project = &project
			}
			err := editAction(os.Stdout, opts)
			if err != nil {
				cli.ShowCommandHelp(ctx, c, "edit")
			}
			return err
		},
	}
}
