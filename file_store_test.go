package todolist

import "testing"

func TestFileStore(t *testing.T) {
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	store.Data[0]
}
