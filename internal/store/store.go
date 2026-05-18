// Package store provides file-backed persistence for quest todos.
//
// Todos are stored as JSON in ~/.quest/todos.json. The package exposes
// helpers to load, save, look up by ID, and remove the store file.
package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrNotFound is returned when a requested todo cannot be found.
var ErrNotFound = errors.New("todo not found")

// Todo is the persisted model for a single task in the quest store.
type Todo struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	Deleted   bool       `json:"deleted"`
	CreatedAt time.Time  `json:"created_at"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	Project   string     `json:"project,omitempty"`
}

func StorePath(subcommand, fileName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".quest", subcommand)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dir, fileName), nil
}

// LoadTodo reads todos from disk, returning an empty slice when the store does not exist.
func LoadTodo() ([]Todo, error) {
	path, err := StorePath("log", "todos.json")
	if err != nil {
		return nil, err
	}

	store, err := os.ReadFile(path)
	if os.IsNotExist(err) {
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

// LoadAndFindIndexByID loads todos and returns the matching todo index for id.
func LoadAndFindIndexByID(id string) ([]Todo, int, error) {
	todos, err := LoadTodo()
	if err != nil {
		return nil, -1, err
	}

	for i := range todos {
		if todos[i].ID == id {
			return todos, i, nil
		}
	}
	return nil, -1, ErrNotFound
}

// SaveTodo writes todos to disk as indented JSON.
func SaveTodo(todos []Todo) error {
	path, err := StorePath("log", "todos.json")
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// NukeTodo removes the todo store file.
func NukeTodo() error {
	path, err := StorePath("log", "todos.json")
	if err != nil {
		return err
	}

	return os.Remove(path)
}
