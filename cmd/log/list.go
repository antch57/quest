package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/antch57/quest/store"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v3"
)

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func isSameDay(a, b time.Time) bool {
	aa := startOfDay(a)
	bb := startOfDay(b)
	return aa.Equal(bb)
}

func printNotebookHeader() {
	accent := text.Colors{text.FgHiCyan}
	doodle := text.Colors{text.FgHiYellow}

	banner := figure.NewFigure("quest log", "weird", true)
	fmt.Print(accent.Sprint(banner.String()))
	wizard := []string{
		"/\\___/\\",
		"( >o.o< )  ~~>~*~*~ FWOOOSH! ~*~*~>~~",
		"\\_____/",
	}
	enemy := []string{
		" .-.",
		"(o_o)",
		"/|_|\\",
	}

	for i := range wizard {
		fmt.Println(doodle.Sprint(fmt.Sprintf("%-39s%s", wizard[i], enemy[i])))
	}
	fmt.Println()
}

func printTable(todos []store.Todo) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault

	t.AppendHeader(table.Row{"ID", "Status", "Title", "Project", "Created", "Due"})

	for _, todo := range todos {
		dueStr := "whenever you want"
		if todo.DueDate != nil {
			dueStr = todo.DueDate.Format("Jan 02, 2006")
		}

		status := "[ ]"
		if todo.Done {
			status = "[x]"
		}

		var color text.Color
		if todo.Done {
			color = text.FgGreen
		} else if todo.DueDate != nil && todo.DueDate.Before(today) {
			color = text.FgRed
		} else {
			color = text.FgYellow
		}

		rowColors := text.Colors{color}
		titleColors := rowColors
		if todo.Done {
			titleColors = text.Colors{color, text.CrossedOut}
		}

		projectStr := todo.Project
		if projectStr == "" {
			projectStr = "-"
		}

		t.AppendRow(table.Row{
			rowColors.Sprint(todo.ID),
			rowColors.Sprint(status),
			titleColors.Sprint(todo.Title),
			rowColors.Sprint(projectStr),
			rowColors.Sprint(todo.CreatedAt.Format("Jan 02, 2006")),
			rowColors.Sprint(dueStr),
		})
	}

	t.Render()
}

func ListCmd() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list all todos",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "show all todos, including done ones",
			},
			&cli.BoolFlag{
				Name:    "done",
				Aliases: []string{"d"},
				Usage:   "show only done todos",
			},
			&cli.BoolFlag{
				Name:  "today",
				Usage: "show todos created today",
			},
			&cli.BoolFlag{
				Name:  "week",
				Usage: "show todos created in the last 7 days",
			},
			&cli.BoolFlag{
				Name:  "month",
				Usage: "show todos created this month",
			},
			&cli.BoolFlag{
				Name:  "overdue",
				Usage: "show todos that are overdue",
			},
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Usage:   "filter todos by project",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			showAll := c.Bool("all")
			showDone := c.Bool("done")
			todayFlag := c.Bool("today")
			weekFlag := c.Bool("week")
			monthFlag := c.Bool("month")
			overdueFlag := c.Bool("overdue")
			projectFilter := c.String("project")

			dateFilterCount := 0
			for _, enabled := range []bool{todayFlag, weekFlag, monthFlag, overdueFlag} {
				if enabled {
					dateFilterCount++
				}
			}
			if dateFilterCount > 1 {
				return fmt.Errorf("use only one date filter at a time: --today, --week, --month, or --overdue")
			}

			now := time.Now()
			today := startOfDay(now)
			weekStart := today.AddDate(0, 0, -6)
			monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
			nextMonthStart := monthStart.AddDate(0, 1, 0)

			matchesDateFilter := func(todo store.Todo) bool {
				if !todayFlag && !weekFlag && !monthFlag && !overdueFlag {
					return true
				}

				created := startOfDay(todo.CreatedAt)

				if todayFlag {
					return isSameDay(created, today)
				}
				if weekFlag {
					return (created.Equal(weekStart) || created.After(weekStart)) && (created.Equal(today) || created.Before(today))
				}
				if monthFlag {
					return (created.Equal(monthStart) || created.After(monthStart)) && created.Before(nextMonthStart)
				}
				if overdueFlag {
					if todo.DueDate == nil {
						return false
					}

					due := startOfDay(*todo.DueDate)
					return !todo.Done && due.Before(today)
				}

				return true
			}

			shouldShow := func(todo store.Todo) bool {
				if todo.Deleted {
					return false
				}
				if !matchesDateFilter(todo) {
					return false
				}
				if projectFilter != "" && todo.Project != projectFilter {
					return false
				}
				if showDone {
					return todo.Done
				}
				if showAll {
					return true
				}
				return !todo.Done
			}

			todos, err := store.Load()
			if err != nil {
				return err
			}

			var filtered []store.Todo
			for _, todo := range todos {
				if shouldShow(todo) {
					filtered = append(filtered, todo)
				}
			}

			if len(filtered) == 0 {
				fmt.Println("No todos found.")
				return nil
			}

			printNotebookHeader()
			printTable(filtered)
			return nil
		},
	}
}
