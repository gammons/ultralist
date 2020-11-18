package ultralist

import (
	"testing"

	"os"
)

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

//SetUpTestMemoryTodoList sets up a fixtures test todolist
func SetUpTestMemoryTodoList() *TodoList {
	store := &MemoryStore{}
	list := &TodoList{}
	list.Data, _ = store.Load()

	todo1 := NewTodo()
	todo1.Subject = "this is the first subject"
	todo1.Projects = []string{"test1"}
	todo1.Contexts = []string{"root"}
	todo1.Due = "2016-04-04"
	todo1.Archive()
	list.Add(todo1)

	todo2 := NewTodo()
	todo2.Subject = "audit userify for 2FA"
	todo2.Projects = []string{"test1"}
	todo2.Contexts = []string{"root", "more"}
	todo2.Complete()
	list.Add(todo2)

	return list
}

func testFileCleanUp() {
	var err = os.Remove(TodosJSONFile)
	if err != nil {
		panic(err)
	}

	return
}
