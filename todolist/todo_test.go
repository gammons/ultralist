package todolist

import "testing"

func TestNewTodo(t *testing.T) {
	todo := NewTodo()

	if todo.Completed || todo.Archived || todo.CompletedDate != "" {
		t.Error("Completed should be false for new todos")
	}
}

func TestValidity(t *testing.T) {
	todo := &Todo{Subject: "test"}
	if !todo.Valid() {
		t.Error("Expected valid todo to be valid")
	}

	invalidTodo := &Todo{Subject: ""}
	if invalidTodo.Valid() {
		t.Error("Invalid todo is being reported as valid")
	}
}
