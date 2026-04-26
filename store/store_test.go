package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// helper function to create a temporary store file
func seedStoreFile(t *testing.T, jsonData string) {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)

	questDir := filepath.Join(home, ".quest")
	if err := os.MkdirAll(questDir, 0o700); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	storeFile := filepath.Join(questDir, "todos.json")
	if err := os.WriteFile(storeFile, []byte(jsonData), 0o600); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}

func TestLoadFileNotExist(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	todos, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != 0 {
		t.Fatalf("expected empty todo list, got %v", todos)
	}
}

func TestLoadSuccess(t *testing.T) {
	seedStoreFile(t, `[{"id":"1","title":"Test Todo","done":false}]`)
	todos, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	if todos[0].ID != "1" || todos[0].Title != "Test Todo" || todos[0].Done != false {
		t.Errorf("unexpected todo data: %v", todos[0])
	}
}

func TestSaveSuccess(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	todos := []Todo{
		{ID: "1", Title: "Test Todo", Done: false},
	}
	if err := Save(todos); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	path, err := storePath()
	if err != nil {
		t.Fatalf("storePath failed: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	var loaded []Todo
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(loaded) != len(todos) {
		t.Fatalf("expected %d todos, got %d", len(todos), len(loaded))
	}
	for i := range todos {
		if loaded[i] != todos[i] {
			t.Errorf("expected todo %v, got %v", todos[i], loaded[i])
		}
	}
}

func TestLoadAndFindIndexByIDSuccess(t *testing.T) {
	seedStoreFile(t, `[{"id":"1","title":"Test Todo","done":false}, {"id":"2","title":"Another Todo","done":true}]`)
	todos, index, err := LoadAndFindIndexByID("1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if index != 0 {
		t.Fatalf("expected index 0, got %d", index)
	}
	if todos[index].ID != "1" || todos[index].Title != "Test Todo" || todos[index].Done != false {
		t.Errorf("unexpected todo data: %v", todos[index])
	}
}

func TestLoadAndFindIndexByIDNotFound(t *testing.T) {
	seedStoreFile(t, `[{"id":"1","title":"Test Todo","done":false}]`)
	todos, index, err := LoadAndFindIndexByID("999")
	if !os.IsNotExist(err) {
		t.Fatalf("expected os.ErrNotExist, got %v", err)
	}
	if index != -1 {
		t.Fatalf("expected index -1 for not found, got %d", index)
	}
	if todos != nil {
		t.Fatalf("expected todos to be nil, got non-nil")
	}
}

func TestSaveAndLoad(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	expected := []Todo{
		{ID: "1", Title: "Test Todo", Done: false},
	}

	if err := Save(expected); err != nil {
		t.Fatalf("failed to save todos: %v", err)
	}

	todos, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(todos) != len(expected) {
		t.Fatalf("expected %d todos, got %d", len(expected), len(todos))
	}
	for i := range expected {
		if todos[i] != expected[i] {
			t.Errorf("expected todo %v, got %v", expected[i], todos[i])
		}
	}
}

func TestNuke(t *testing.T) {
	seedStoreFile(t, `[{"id":"1","title":"Test Todo","done":false}]`)
	if err := Nuke(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	path, err := storePath()
	if err != nil {
		t.Fatalf("storePath failed: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed, but it exists")
	}
}
