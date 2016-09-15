package todolist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStore(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	todos, _ := store.Load()
	assert.Equal(todos[0].Subject, "this is the first subject", "")
}

func TestSave(t *testing.T) {
	store := &FileStore{FileLocation: "todos.json"}
	todos, _ := store.Load()
	store.Save(todos)
}
