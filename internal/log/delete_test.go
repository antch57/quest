package log

import (
	"os"
	"strings"
	"testing"

	"github.com/antch57/quest/internal/store"
)

func Test_deleteAction(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name      string
		args      args
		seedTodos []store.Todo
		wantErr   bool
		wantTodo  store.Todo
	}{
		{
			name: "delete existing task",
			args: args{
				id: "1",
			},
			seedTodos: []store.Todo{
				{ID: "1", Title: "existing task"},
				{ID: "2", Title: "other task"},
			},
			wantErr:  false,
			wantTodo: store.Todo{ID: "1", Title: "existing task", Deleted: true},
		},
		{
			name: "delete non-existing task",
			args: args{
				id: "99",
			},
			seedTodos: []store.Todo{
				{ID: "1", Title: "existing task"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			seedStoreTodos(t, tt.seedTodos)

			if err := deleteAction(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("deleteAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			todos, idx, err := store.LoadAndFindIndexByID(tt.args.id)
			if err != nil {
				t.Fatalf("store.LoadAndFindIndexByID() error = %v", err)
			}
			if got := todos[idx]; got.ID != tt.wantTodo.ID || got.Title != tt.wantTodo.Title || got.Deleted != tt.wantTodo.Deleted {
				t.Errorf("deleted todo = %+v, want %+v", got, tt.wantTodo)
			}
		})
	}
}

func Test_nukeAction(t *testing.T) {
	type args struct {
		response string
	}
	tests := []struct {
		name           string
		args           args
		seedTodos      []store.Todo
		wantErr        bool
		wantStoreExist bool
	}{
		{
			name: "nuke all tasks",
			args: args{
				response: "y\n",
			},
			seedTodos:      []store.Todo{{ID: "1", Title: "existing task"}},
			wantErr:        false,
			wantStoreExist: false,
		},
		{
			name: "abort nuke",
			args: args{
				response: "n\n",
			},
			seedTodos:      []store.Todo{{ID: "1", Title: "existing task"}},
			wantErr:        false,
			wantStoreExist: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			seedStoreTodos(t, tt.seedTodos)

			if err := nukeAction(strings.NewReader(tt.args.response)); (err != nil) != tt.wantErr {
				t.Errorf("nukeAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			path := os.Getenv("HOME") + "/.quest/todos.json"

			_, err := os.Stat(path)
			gotStoreExist := !os.IsNotExist(err)
			if err != nil && !os.IsNotExist(err) {
				t.Fatalf("os.Stat() error = %v", err)
			}
			if gotStoreExist != tt.wantStoreExist {
				t.Errorf("store exists = %v, want %v", gotStoreExist, tt.wantStoreExist)
			}
		})
	}
}
