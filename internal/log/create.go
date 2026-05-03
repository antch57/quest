package log

import (
	"context"
	"fmt"
	"time"

	"github.com/antch57/quest/internal/store"
	"github.com/urfave/cli/v3"
)

type CreateOptions struct {
	Title   string
	Due     string
	Project string
}

func createID(todos []store.Todo) string {
	max := 0
	for _, todo := range todos {
		var idInt int
		_, err := fmt.Sscanf(todo.ID, "%d", &idInt)
		if err == nil && idInt > max {
			max = idInt
		}
	}
	return fmt.Sprintf("%d", max+1)
}

func createTodo(opts CreateOptions) error {
	if opts.Title == "" {
		return fmt.Errorf("title is required")
	}

	var dueDate *time.Time
	if opts.Due != "" {
		parsedDueDate, err := time.Parse("01-02-2006", opts.Due)
		if err != nil {
			return fmt.Errorf("invalid due date format (expected mm-dd-yyyy): %v", err)
		}
		dueDate = &parsedDueDate
	}

	todos, err := store.Load()
	if err != nil {
		return err
	}

	todos = append(todos, store.Todo{
		ID:        createID(todos),
		Title:     opts.Title,
		CreatedAt: time.Now(),
		DueDate:   dueDate,
		Project:   opts.Project,
	})

	if err := store.Save(todos); err != nil {
		return err
	}

	fmt.Printf("created todo: %s\n", opts.Title)
	return nil
}

func CreateCmd() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "create a new todo item.",
		UsageText: `quest log create --title "buy groceries"`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "title",
				Aliases:  []string{"t"},
				Usage:    "title for the todo (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "due",
				Aliases: []string{"d"},
				Usage:   "due date for the todo (format: mm-dd-yyyy)",
			},
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Usage:   "project or folder for this todo",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			opts := CreateOptions{
				Title:   c.String("title"),
				Project: c.String("project"),
				Due:     c.String("due"),
			}
			err := createTodo(opts)
			if err != nil {
				cli.ShowCommandHelp(ctx, c, "create")
			}
			return err
		},
	}
}
