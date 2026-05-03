package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/antch57/quest/internal/store"
	"github.com/common-nighthawk/go-figure"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/urfave/cli/v3"
)

type ListOptions struct {
	ShowAll       bool
	ShowDone      bool
	Today         bool
	Week          bool
	Month         bool
	Overdue       bool
	ProjectFilter string
}

func ListCmd() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "list all todos",
		UsageText: "quest log list [options]",
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
			opts := ListOptions{
				ShowAll:       c.Bool("all"),
				ShowDone:      c.Bool("done"),
				Today:         c.Bool("today"),
				Week:          c.Bool("week"),
				Month:         c.Bool("month"),
				Overdue:       c.Bool("overdue"),
				ProjectFilter: c.String("project"),
			}
			err := listAction(os.Stdout, opts)
			if err != nil {
				cli.ShowCommandHelp(ctx, c, "list")
			}
			return err
		},
	}
}

func listAction(w io.Writer, opts ListOptions) error {
	if err := validateDateFilter(opts); err != nil {
		return err
	}
	todos, err := store.Load()
	if err != nil {
		return err
	}

	filtered := filterTodos(todos, opts)
	if len(filtered) == 0 {
		fmt.Fprintln(w, "no todos found")
		return nil
	}
	printNotebookHeader(w)
	printTable(w, filtered)
	return nil
}

func filterTodos(todos []store.Todo, opts ListOptions) []store.Todo {
	var filtered []store.Todo
	for _, todo := range todos {
		if shouldShow(opts, todo) {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

func validateDateFilter(opts ListOptions) error {
	count := 0
	for _, enabled := range []bool{opts.Today, opts.Week, opts.Month, opts.Overdue} {
		if enabled {
			count++
		}
	}
	if count > 1 {
		return fmt.Errorf("use only one date filter at a time: --today, --week, --month, or --overdue")
	}
	return nil
}

func shouldShow(opts ListOptions, todo store.Todo) bool {
	if todo.Deleted {
		return false
	}
	if !matchesDateFilter(opts, todo) {
		return false
	}
	if opts.ProjectFilter != "" && todo.Project != opts.ProjectFilter {
		return false
	}
	if opts.ShowDone {
		return todo.Done
	}
	if opts.ShowAll {
		return true
	}
	return !todo.Done
}

func printNotebookHeader(w io.Writer) {
	accent := text.Colors{text.FgHiCyan}
	doodle := text.Colors{text.FgHiYellow}

	banner := figure.NewFigure("quest log", "weird", true)
	fmt.Fprint(w, accent.Sprint(banner.String()))
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
		fmt.Fprintln(w, doodle.Sprint(fmt.Sprintf("%-39s%s", wizard[i], enemy[i])))
	}
	fmt.Fprintln(w)
}

func printTable(w io.Writer, todos []store.Todo) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault

	t.AppendHeader(table.Row{"id", "status", "title", "project", "created", "due"})

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

func matchesDateFilter(opts ListOptions, todo store.Todo) bool {
	now := time.Now()
	today := startOfDay(now)
	weekStart := today.AddDate(0, 0, -6)
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	if !opts.Today && !opts.Week && !opts.Month && !opts.Overdue {
		return true
	}

	created := startOfDay(todo.CreatedAt)

	if opts.Today {
		return isSameDay(created, today)
	}
	if opts.Week {
		return (created.Equal(weekStart) || created.After(weekStart)) && (created.Equal(today) || created.Before(today))
	}
	if opts.Month {
		return (created.Equal(monthStart) || created.After(monthStart)) && created.Before(nextMonthStart)
	}
	if opts.Overdue {
		if todo.DueDate == nil {
			return false
		}

		due := startOfDay(*todo.DueDate)
		return !todo.Done && due.Before(today)
	}

	return true
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func isSameDay(a, b time.Time) bool {
	aa := startOfDay(a)
	bb := startOfDay(b)
	return aa.Equal(bb)
}
