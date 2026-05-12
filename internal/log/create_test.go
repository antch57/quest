package log

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/antch57/quest/internal/store"
)

func useTempHome(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
}

func seedStoreJSON(t *testing.T, jsonData string) {
	t.Helper()
	home := os.Getenv("HOME")
	if home == "" {
		t.Fatal("HOME is not set")
	}
	questDir := filepath.Join(home, ".quest")
	if err := os.MkdirAll(questDir, 0o700); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	path := filepath.Join(questDir, "todos.json")
	if err := os.WriteFile(path, []byte(jsonData), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
}

func seedStoreTodos(t *testing.T, todos []store.Todo) {
	t.Helper()
	data, err := json.Marshal(todos)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	seedStoreJSON(t, string(data))
}

func Test_createID(t *testing.T) {
	type args struct {
		todos []store.Todo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid input",
			args: args{
				todos: []store.Todo{{ID: "123", Title: "Task 1"}, {ID: "456", Title: "Task 2"}},
			},
			want: "457",
		},
		{
			name: "edge case zero",
			args: args{
				todos: []store.Todo{},
			},
			want: "1",
		},
		{
			name: "edge case empty",
			args: args{
				todos: []store.Todo{{ID: "", Title: "Task 1"}},
			},
			want: "1",
		},
		{
			name: "edge case nil",
			args: args{
				todos: []store.Todo{{ID: "456", Title: "Task 2"}, {ID: "", Title: "Task 3"}},
			},
			want: "457",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createID(tt.args.todos); got != tt.want {
				t.Errorf("createID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createTodo(t *testing.T) {
	type args struct {
		opts CreateOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				opts: CreateOptions{Title: "Test Todo", Due: "01-02-2006"},
			},
			wantErr: false,
		},
		{
			name: "invalid due date format",
			args: args{
				opts: CreateOptions{Title: "Test Todo", Due: "invalid_due_date_format"},
			},
			wantErr: true,
		},
		{
			name: "empty title",
			args: args{
				opts: CreateOptions{Title: "", Due: "01-02-2006"},
			},
			wantErr: true,
		},
		{
			name: "empty due date",
			args: args{
				opts: CreateOptions{Title: "Test Todo", Due: ""},
			},
			wantErr: false,
		},
		{
			name: "empty project",
			args: args{
				opts: CreateOptions{Title: "Test Todo", Project: ""},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)

			if err := createTodo(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("createTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
