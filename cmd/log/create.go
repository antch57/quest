package log

import (
	"context"
	"fmt"
	"time"

	"github.com/antch57/quest/store"
	"github.com/urfave/cli/v3"
)

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
			title := c.String("title")
			project := c.String("project")
			dueDateStr := c.String("due")

			var dueDate *time.Time
			if dueDateStr != "" {
				parsedDueDate, err := time.Parse("01-02-2006", dueDateStr)
				if err != nil {
					cli.ShowCommandHelp(ctx, c, "create")
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
				Title:     title,
				CreatedAt: time.Now(),
				DueDate:   dueDate,
				Project:   project,
			})
			if err := store.Save(todos); err != nil {
				return err
			}
			fmt.Printf("created todo: %s\n", title)
			return nil
		},
	}
}
