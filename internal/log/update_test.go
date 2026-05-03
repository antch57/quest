package log

import (
	"io"
	"testing"
	"time"

	"github.com/antch57/quest/internal/store"
)

func Test_doneAction(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name      string
		args      args
		seedTodos []store.Todo
		wantErr   bool
		wantDone  bool
	}{
		{
			name:      "mark existing task done",
			args:      args{id: "1"},
			seedTodos: []store.Todo{{ID: "1", Title: "task", Done: false, CreatedAt: time.Now()}},
			wantErr:   false,
			wantDone:  true,
		},
		{
			name:      "mark already done task",
			args:      args{id: "1"},
			seedTodos: []store.Todo{{ID: "1", Title: "task", Done: true, CreatedAt: time.Now()}},
			wantErr:   false,
			wantDone:  true,
		},
		{
			name:      "missing task id returns error",
			args:      args{id: "99"},
			seedTodos: []store.Todo{{ID: "1", Title: "task", Done: false, CreatedAt: time.Now()}},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			seedStoreTodos(t, tt.seedTodos)

			if err := doneAction(io.Discard, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("doneAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			todos, idx, err := store.LoadAndFindIndexByID(tt.args.id)
			if err != nil {
				t.Fatalf("store.LoadAndFindIndexByID() error = %v", err)
			}
			if got := todos[idx].Done; got != tt.wantDone {
				t.Errorf("todo.Done = %v, want %v", got, tt.wantDone)
			}
		})
	}
}

func Test_editAction(t *testing.T) {
	baseDue := time.Date(2026, time.May, 10, 0, 0, 0, 0, time.UTC)

	type args struct {
		opts EditOptions
	}
	tests := []struct {
		name       string
		args       args
		seedTodos  []store.Todo
		wantErr    bool
		wantTitle  string
		wantDone   bool
		wantDue    string
		wantNilDue bool
		wantProj   string
	}{
		{
			name:      "edit missing id returns error",
			args:      args{opts: EditOptions{ID: "99"}},
			seedTodos: []store.Todo{{ID: "1", Title: "task", CreatedAt: time.Now()}},
			wantErr:   true,
		},
		{
			name:      "done and undone flags conflict",
			args:      args{opts: EditOptions{ID: "1", Done: true, Undone: true}},
			seedTodos: []store.Todo{{ID: "1", Title: "task", CreatedAt: time.Now()}},
			wantErr:   true,
		},
		{
			name:      "clear due and due conflict",
			args:      args{opts: EditOptions{ID: "1", ClearDue: true, Due: strPtr("05-20-2026")}},
			seedTodos: []store.Todo{{ID: "1", Title: "task", DueDate: &baseDue, CreatedAt: time.Now()}},
			wantErr:   true,
		},
		{
			name:      "invalid due format",
			args:      args{opts: EditOptions{ID: "1", Due: strPtr("bad-date")}},
			seedTodos: []store.Todo{{ID: "1", Title: "task", CreatedAt: time.Now()}},
			wantErr:   true,
		},
		{
			name:       "update title",
			args:       args{opts: EditOptions{ID: "1", Title: strPtr("new title")}},
			seedTodos:  []store.Todo{{ID: "1", Title: "old title", CreatedAt: time.Now()}},
			wantErr:    false,
			wantTitle:  "new title",
			wantDone:   false,
			wantNilDue: true,
		},
		{
			name:      "set due date",
			args:      args{opts: EditOptions{ID: "1", Due: strPtr("05-20-2026")}},
			seedTodos: []store.Todo{{ID: "1", Title: "task", CreatedAt: time.Now()}},
			wantErr:   false,
			wantTitle: "task",
			wantDone:  false,
			wantDue:   "05-20-2026",
		},
		{
			name:       "clear due date",
			args:       args{opts: EditOptions{ID: "1", ClearDue: true}},
			seedTodos:  []store.Todo{{ID: "1", Title: "task", DueDate: &baseDue, CreatedAt: time.Now()}},
			wantErr:    false,
			wantTitle:  "task",
			wantDone:   false,
			wantNilDue: true,
		},
		{
			name:       "mark done",
			args:       args{opts: EditOptions{ID: "1", Done: true}},
			seedTodos:  []store.Todo{{ID: "1", Title: "task", Done: false, CreatedAt: time.Now()}},
			wantErr:    false,
			wantTitle:  "task",
			wantDone:   true,
			wantNilDue: true,
		},
		{
			name:       "mark undone",
			args:       args{opts: EditOptions{ID: "1", Undone: true}},
			seedTodos:  []store.Todo{{ID: "1", Title: "task", Done: true, CreatedAt: time.Now()}},
			wantErr:    false,
			wantTitle:  "task",
			wantDone:   false,
			wantNilDue: true,
		},
		{
			name:       "update project",
			args:       args{opts: EditOptions{ID: "1", Project: strPtr("work")}},
			seedTodos:  []store.Todo{{ID: "1", Title: "task", Project: "home", CreatedAt: time.Now()}},
			wantErr:    false,
			wantTitle:  "task",
			wantDone:   false,
			wantNilDue: true,
			wantProj:   "work",
		},
		{
			name:      "combined edit fields",
			args:      args{opts: EditOptions{ID: "1", Title: strPtr("combined"), Project: strPtr("ops"), Done: true, Due: strPtr("05-20-2026")}},
			seedTodos: []store.Todo{{ID: "1", Title: "task", Done: false, CreatedAt: time.Now()}},
			wantErr:   false,
			wantTitle: "combined",
			wantDone:  true,
			wantDue:   "05-20-2026",
			wantProj:  "ops",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			seedStoreTodos(t, tt.seedTodos)

			if err := editAction(io.Discard, tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("editAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			todos, idx, err := store.LoadAndFindIndexByID(tt.args.opts.ID)
			if err != nil {
				t.Fatalf("store.LoadAndFindIndexByID() error = %v", err)
			}
			got := todos[idx]

			if tt.wantTitle != "" && got.Title != tt.wantTitle {
				t.Errorf("todo.Title = %q, want %q", got.Title, tt.wantTitle)
			}
			if got.Done != tt.wantDone {
				t.Errorf("todo.Done = %v, want %v", got.Done, tt.wantDone)
			}
			if tt.wantProj != "" && got.Project != tt.wantProj {
				t.Errorf("todo.Project = %q, want %q", got.Project, tt.wantProj)
			}
			if tt.wantNilDue && got.DueDate != nil {
				t.Errorf("todo.DueDate = %v, want nil", got.DueDate)
			}
			if tt.wantDue != "" {
				if got.DueDate == nil {
					t.Fatalf("todo.DueDate = nil, want %q", tt.wantDue)
				}
				if gotDate := got.DueDate.Format("01-02-2006"); gotDate != tt.wantDue {
					t.Errorf("todo.DueDate = %q, want %q", gotDate, tt.wantDue)
				}
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
