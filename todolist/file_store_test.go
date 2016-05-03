package todolist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStore(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	assert.Equal(store.Data[0].Subject, "this is the first subject", "")
}

func TestSave(t *testing.T) {
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	store.Save()
}

func TestNextId(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	assert.Equal(3, store.NextId())
}

func TestIndexOf(t *testing.T) {
	assert := assert.New(t)
	todo := &Todo{Subject: "Grant"}
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()

	assert.Equal(-1, store.IndexOf(todo))
	assert.Equal(0, store.IndexOf(store.Data[0]))
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	assert.Equal(2, len(store.Data))
	store.Delete(1)
	assert.Equal(1, len(store.Data))
}

func TestComplete(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	assert.Equal(false, store.FindById(1).Completed)
	store.Complete(1)
	assert.Equal(true, store.FindById(1).Completed)
}

func TestArchive(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	assert.Equal(false, store.FindById(2).Archived)
	store.Archive(2)
	assert.Equal(true, store.FindById(2).Archived)
}
func TestUnarchive(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	assert.Equal(true, store.FindById(1).Archived)
	store.Unarchive(1)
	assert.Equal(false, store.FindById(1).Archived)
}

func TestUncomplete(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	assert.Equal(true, store.FindById(2).Completed)
	store.Uncomplete(2)
	assert.Equal(false, store.FindById(2).Completed)
}
