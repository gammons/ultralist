package todolist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextId(t *testing.T) {
	assert := assert.New(t)
	list := &TodoList{}
	assert.Equal(1, list.NextId())
}

func TestIndexOf(t *testing.T) {
	assert := assert.New(t)
	todo := &Todo{Subject: "Grant"}
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	todos, _ := store.Load()
	list.Load(todos)

	assert.Equal(-1, list.IndexOf(todo))
	assert.Equal(0, list.IndexOf(list.Data[0]))
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	todos, _ := store.Load()
	list.Load(todos)
	assert.Equal(2, len(list.Data))
	list.Delete(1)
	assert.Equal(1, len(list.Data))
}

func TestComplete(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	todos, _ := store.Load()
	list.Load(todos)
	assert.Equal(false, list.FindById(1).Completed)
	list.Complete(1)
	assert.Equal(true, list.FindById(1).Completed)
}

func TestArchive(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	todos, _ := store.Load()
	list.Load(todos)
	assert.Equal(false, list.FindById(2).Archived)
	list.Archive(2)
	assert.Equal(true, list.FindById(2).Archived)
}
func TestUnarchive(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	todos, _ := store.Load()
	list.Load(todos)
	assert.Equal(true, list.FindById(1).Archived)
	list.Unarchive(1)
	assert.Equal(false, list.FindById(1).Archived)
}

func TestUncomplete(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	todos, _ := store.Load()
	list.Load(todos)
	assert.Equal(true, list.FindById(2).Completed)
	list.Uncomplete(2)
	assert.Equal(false, list.FindById(2).Completed)
}
