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
	assert.Equal(0, store.IndexOf(&store.Data[0]))
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	store.Load()
	assert.Equal(2, len(store.Data))
	store.Delete(1)
	assert.Equal(1, len(store.Data))
}
