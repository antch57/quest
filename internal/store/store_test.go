package store

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
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

func timePtr(t time.Time) *time.Time {
	return &t
}

func Test_storePath(t *testing.T) {
	tests := []struct {
		name       string
		wantSuffix string
		wantErr    bool
	}{
		{
			name:       "returns todos path under quest dir",
			wantSuffix: filepath.Join(".quest", "todos.json"),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)

			got, err := storePath()
			if (err != nil) != tt.wantErr {
				t.Fatalf("storePath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			home := os.Getenv("HOME")
			want := filepath.Join(home, tt.wantSuffix)
			if got != want {
				t.Errorf("storePath() = %v, want %v", got, want)
			}

			if info, err := os.Stat(filepath.Dir(got)); err != nil {
				t.Fatalf("os.Stat() error = %v", err)
			} else if !info.IsDir() {
				t.Fatalf("parent directory was not created")
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		seed    string
		want    []Todo
		wantErr bool
	}{
		{
			name:    "missing file returns empty todos",
			want:    []Todo{},
			wantErr: false,
		},
		{
			name: "valid json returns todos",
			seed: `[
				{"id":"1","title":"task one","done":false},
				{"id":"2","title":"task two","done":true}
			]`,
			want: []Todo{
				{ID: "1", Title: "task one", Done: false},
				{ID: "2", Title: "task two", Done: true},
			},
			wantErr: false,
		},
		{
			name:    "corrupt json returns error",
			seed:    "{broken-json}",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			if tt.seed != "" {
				seedStoreJSON(t, tt.seed)
			}

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadAndFindIndexByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		seed       string
		args       args
		want       []Todo
		want1      int
		wantErr    bool
		wantNoFile bool
	}{
		{
			name: "finds existing id",
			seed: `[
				{"id":"1","title":"task one","done":false},
				{"id":"2","title":"task two","done":true}
			]`,
			args: args{id: "2"},
			want: []Todo{
				{ID: "1", Title: "task one", Done: false},
				{ID: "2", Title: "task two", Done: true},
			},
			want1:   1,
			wantErr: false,
		},
		{
			name: "missing id returns ErrNotFound",
			seed: `[
				{"id":"1","title":"task one","done":false}
			]`,
			args:       args{id: "999"},
			want1:      -1,
			wantErr:    true,
			wantNoFile: true,
		},
		{
			name:       "corrupt json bubbles load error",
			seed:       "{broken-json}",
			args:       args{id: "1"},
			want1:      -1,
			wantErr:    true,
			wantNoFile: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			if tt.seed != "" {
				seedStoreJSON(t, tt.seed)
			}

			got, got1, err := LoadAndFindIndexByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadAndFindIndexByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.wantNoFile && !errors.Is(err, ErrNotFound) {
				t.Fatalf("LoadAndFindIndexByID() error = %v, want ErrNotFound", err)
			}
			if tt.wantErr {
				if got1 != tt.want1 {
					t.Errorf("LoadAndFindIndexByID() got1 = %v, want %v", got1, tt.want1)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAndFindIndexByID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LoadAndFindIndexByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSave(t *testing.T) {
	type args struct {
		todos []Todo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "save empty todos",
			args:    args{todos: []Todo{}},
			wantErr: false,
		},
		{
			name: "save and verify round-trip",
			args: args{todos: []Todo{
				{
					ID:        "1",
					Title:     "task one",
					Done:      false,
					Deleted:   false,
					CreatedAt: time.Date(2026, time.May, 2, 12, 0, 0, 0, time.UTC),
					DueDate:   timePtr(time.Date(2026, time.May, 10, 0, 0, 0, 0, time.UTC)),
					Project:   "work",
				},
				{
					ID:        "2",
					Title:     "task two",
					Done:      true,
					Deleted:   false,
					CreatedAt: time.Date(2026, time.May, 3, 12, 0, 0, 0, time.UTC),
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)

			if err := Save(tt.args.todos); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			got, err := Load()
			if err != nil {
				t.Fatalf("Load() after Save() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.args.todos) {
				t.Errorf("Save()/Load() got = %v, want %v", got, tt.args.todos)
			}
		})
	}
}

func TestNuke(t *testing.T) {
	tests := []struct {
		name          string
		seed          string
		wantErr       bool
		wantNotExists bool
	}{
		{
			name:          "nukes existing store file",
			seed:          `[{"id":"1","title":"task one","done":false}]`,
			wantErr:       false,
			wantNotExists: true,
		},
		{
			name:    "returns error when store file does not exist",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempHome(t)
			if tt.seed != "" {
				seedStoreJSON(t, tt.seed)
			}

			if err := Nuke(); (err != nil) != tt.wantErr {
				t.Errorf("Nuke() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantNotExists {
				path, err := storePath()
				if err != nil {
					t.Fatalf("storePath() error = %v", err)
				}
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Fatalf("expected store file to be removed, stat err = %v", err)
				}
			}
		})
	}
}
