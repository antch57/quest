package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Todo struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	Deleted   bool       `json:"deleted"`
	CreatedAt time.Time  `json:"created_at"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	Project   string     `json:"project,omitempty"`
}

func storePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".quest")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dir, "todos.json"), nil
}

func Load() ([]Todo, error) {
	path, err := storePath()
	if err != nil {
		return nil, err
	}

	store, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// If the file does not exist, return an empty list of todos
		return []Todo{}, nil
	}
	if err != nil {
		return nil, err
	}

	var todos []Todo
	if err := json.Unmarshal(store, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func LoadAndFindIndexByID(id string) ([]Todo, int, error) {
	todos, err := Load()
	if err != nil {
		return nil, -1, err
	}

	for i := range todos {
		if todos[i].ID == id {
			return todos, i, nil
		}
	}

	return nil, -1, os.ErrNotExist
}

func Save(todos []Todo) error {
	path, err := storePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func Nuke() error {
	path, err := storePath()
	if err != nil {
		return err
	}

	return os.Remove(path)
}
