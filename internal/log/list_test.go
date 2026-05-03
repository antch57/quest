package log

import (
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/antch57/quest/internal/store"
)

func Test_listAction(t *testing.T) {
	type args struct {
		opts ListOptions
	}
	tests := []struct {
		name      string
		args      args
		seedTodos []store.Todo
		wantErr   bool
	}{
		{
			name: "list all tasks",
			args: args{opts: ListOptions{}},
			seedTodos: []store.Todo{
				{ID: "1", Title: "todo", CreatedAt: time.Now()},
				{ID: "2", Title: "done", Done: true, CreatedAt: time.Now()},
				{ID: "3", Title: "deleted", Deleted: true, CreatedAt: time.Now()},
			},
			wantErr: false,
		},
		{
			name: "list tasks with today filter",
			args: args{opts: ListOptions{Today: true}},
			seedTodos: []store.Todo{
				{ID: "1", Title: "today task", CreatedAt: time.Now()},
				{ID: "2", Title: "older task", CreatedAt: time.Now().AddDate(0, 0, -10)},
			},
			wantErr: false,
		},
		{
			name:      "list tasks with conflicting date filters",
			args:      args{opts: ListOptions{Today: true, Week: true}},
			seedTodos: []store.Todo{{ID: "1", Title: "task", CreatedAt: time.Now()}},
			wantErr:   true,
		},
		{
			name: "list tasks with done filter",
			args: args{opts: ListOptions{ShowDone: true}},
			seedTodos: []store.Todo{
				{ID: "1", Title: "todo", Done: false, CreatedAt: time.Now()},
				{ID: "2", Title: "done", Done: true, CreatedAt: time.Now()},
			},
			wantErr: false,
		},
		{
			name: "list tasks with project filter",
			args: args{opts: ListOptions{ProjectFilter: "work"}},
			seedTodos: []store.Todo{
				{ID: "1", Title: "work item", Project: "work", CreatedAt: time.Now()},
				{ID: "2", Title: "home item", Project: "home", CreatedAt: time.Now()},
			},
			wantErr: false,
		},
		{
			name:      "empty store prints no todos found",
			args:      args{opts: ListOptions{}},
			seedTodos: []store.Todo{},
			wantErr:   false,
		},
		{
			name: "all todos filtered out prints no todos found",
			args: args{opts: ListOptions{}},
			seedTodos: []store.Todo{
				{ID: "1", Title: "deleted task", Deleted: true, CreatedAt: time.Now()},
				{ID: "2", Title: "done task", Done: true, CreatedAt: time.Now()},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			seedStoreTodos(t, tt.seedTodos)

			if err := listAction(io.Discard, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("listAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_filterTodos(t *testing.T) {
	type args struct {
		todos []store.Todo
		opts  ListOptions
	}
	now := time.Now()
	todos := []store.Todo{
		{ID: "1", Title: "todo", Done: false, Deleted: false, Project: "work", CreatedAt: now},
		{ID: "2", Title: "done", Done: true, Deleted: false, Project: "home", CreatedAt: now},
		{ID: "3", Title: "deleted", Done: false, Deleted: true, Project: "work", CreatedAt: now},
	}
	tests := []struct {
		name string
		args args
		want []store.Todo
	}{
		{
			name: "default shows only active not-done todos",
			args: args{todos: todos, opts: ListOptions{}},
			want: []store.Todo{todos[0]},
		},
		{
			name: "show done returns only done todos",
			args: args{todos: todos, opts: ListOptions{ShowDone: true}},
			want: []store.Todo{todos[1]},
		},
		{
			name: "show all returns all non-deleted todos",
			args: args{todos: todos, opts: ListOptions{ShowAll: true}},
			want: []store.Todo{todos[0], todos[1]},
		},
		{
			name: "project filter returns matching non-deleted todos",
			args: args{todos: todos, opts: ListOptions{ShowAll: true, ProjectFilter: "work"}},
			want: []store.Todo{todos[0]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterTodos(tt.args.todos, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterTodos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateDateFilter(t *testing.T) {
	type args struct {
		opts ListOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no date filter",
			args:    args{opts: ListOptions{}},
			wantErr: false,
		},
		{
			name:    "single date filter",
			args:    args{opts: ListOptions{Today: true}},
			wantErr: false,
		},
		{
			name:    "multiple date filters returns error",
			args:    args{opts: ListOptions{Today: true, Week: true}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDateFilter(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("validateDateFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_shouldShow(t *testing.T) {
	type args struct {
		opts ListOptions
		todo store.Todo
	}
	now := time.Now()
	overdue := now.AddDate(0, 0, -1)
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "deleted todo is hidden",
			args: args{opts: ListOptions{}, todo: store.Todo{Deleted: true, CreatedAt: now}},
			want: false,
		},
		{
			name: "show done includes done todo",
			args: args{opts: ListOptions{ShowDone: true}, todo: store.Todo{Done: true, CreatedAt: now}},
			want: true,
		},
		{
			name: "show done excludes not-done todo",
			args: args{opts: ListOptions{ShowDone: true}, todo: store.Todo{Done: false, CreatedAt: now}},
			want: false,
		},
		{
			name: "show all includes active todo",
			args: args{opts: ListOptions{ShowAll: true}, todo: store.Todo{Done: false, CreatedAt: now}},
			want: true,
		},
		{
			name: "default hides done todo",
			args: args{opts: ListOptions{}, todo: store.Todo{Done: true, CreatedAt: now}},
			want: false,
		},
		{
			name: "project mismatch is hidden",
			args: args{opts: ListOptions{ProjectFilter: "work"}, todo: store.Todo{Project: "home", CreatedAt: now}},
			want: false,
		},
		{
			name: "overdue filter matches overdue active todo",
			args: args{opts: ListOptions{Overdue: true}, todo: store.Todo{DueDate: &overdue, Done: false, CreatedAt: now}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldShow(tt.args.opts, tt.args.todo); got != tt.want {
				t.Errorf("shouldShow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printNotebookHeader(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "prints banner without panic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printNotebookHeader(io.Discard)
		})
	}
}

func Test_printTable(t *testing.T) {
	type args struct {
		todos []store.Todo
	}
	now := time.Now()
	due := now.AddDate(0, 0, 2)
	tests := []struct {
		name string
		args args
	}{
		{
			name: "prints table with rows",
			args: args{todos: []store.Todo{{ID: "1", Title: "task", CreatedAt: now, DueDate: &due}}},
		},
		{
			name: "prints table when empty",
			args: args{todos: []store.Todo{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printTable(io.Discard, tt.args.todos)
		})
	}
}

func Test_matchesDateFilter(t *testing.T) {
	type args struct {
		opts ListOptions
		todo store.Todo
	}
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	threeDaysAgo := now.AddDate(0, 0, -3)
	eightDaysAgo := now.AddDate(0, 0, -8)
	monthStart := time.Date(now.Year(), now.Month(), 1, 10, 0, 0, 0, now.Location())
	lastMonth := monthStart.AddDate(0, -1, 0)
	tomorrow := now.AddDate(0, 0, 1)
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no filter returns true",
			args: args{opts: ListOptions{}, todo: store.Todo{CreatedAt: eightDaysAgo}},
			want: true,
		},
		{
			name: "today filter matches today",
			args: args{opts: ListOptions{Today: true}, todo: store.Todo{CreatedAt: now}},
			want: true,
		},
		{
			name: "today filter excludes non-today",
			args: args{opts: ListOptions{Today: true}, todo: store.Todo{CreatedAt: yesterday}},
			want: false,
		},
		{
			name: "week filter includes recent",
			args: args{opts: ListOptions{Week: true}, todo: store.Todo{CreatedAt: threeDaysAgo}},
			want: true,
		},
		{
			name: "week filter excludes older than 7 days",
			args: args{opts: ListOptions{Week: true}, todo: store.Todo{CreatedAt: eightDaysAgo}},
			want: false,
		},
		{
			name: "month filter includes this month",
			args: args{opts: ListOptions{Month: true}, todo: store.Todo{CreatedAt: monthStart}},
			want: true,
		},
		{
			name: "month filter excludes previous month",
			args: args{opts: ListOptions{Month: true}, todo: store.Todo{CreatedAt: lastMonth}},
			want: false,
		},
		{
			name: "overdue matches overdue not-done todo",
			args: args{opts: ListOptions{Overdue: true}, todo: store.Todo{DueDate: &yesterday, Done: false, CreatedAt: now}},
			want: true,
		},
		{
			name: "overdue excludes done todo",
			args: args{opts: ListOptions{Overdue: true}, todo: store.Todo{DueDate: &yesterday, Done: true, CreatedAt: now}},
			want: false,
		},
		{
			name: "overdue excludes missing due date",
			args: args{opts: ListOptions{Overdue: true}, todo: store.Todo{DueDate: nil, Done: false, CreatedAt: now}},
			want: false,
		},
		{
			name: "overdue excludes future due date",
			args: args{opts: ListOptions{Overdue: true}, todo: store.Todo{DueDate: &tomorrow, Done: false, CreatedAt: now}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesDateFilter(tt.args.opts, tt.args.todo); got != tt.want {
				t.Errorf("matchesDateFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_startOfDay(t *testing.T) {
	type args struct {
		t time.Time
	}
	loc := time.FixedZone("UTC-7", -7*60*60)
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "resets clock to midnight",
			args: args{t: time.Date(2026, time.May, 2, 14, 35, 12, 123, loc)},
			want: time.Date(2026, time.May, 2, 0, 0, 0, 0, loc),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := startOfDay(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("startOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSameDay(t *testing.T) {
	type args struct {
		a time.Time
		b time.Time
	}
	loc := time.FixedZone("UTC-7", -7*60*60)
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same date different times",
			args: args{
				a: time.Date(2026, time.May, 2, 1, 0, 0, 0, loc),
				b: time.Date(2026, time.May, 2, 23, 59, 0, 0, loc),
			},
			want: true,
		},
		{
			name: "different dates",
			args: args{
				a: time.Date(2026, time.May, 2, 23, 59, 0, 0, loc),
				b: time.Date(2026, time.May, 3, 0, 0, 0, 0, loc),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSameDay(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("isSameDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
