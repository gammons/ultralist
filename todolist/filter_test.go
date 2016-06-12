package todolist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterArchived(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	list.Load(store.Load())
	filter := NewFilter(list.Todos())
	archived := filter.filterArchived("l archived")
	assert.Equal(1, len(archived))
	assert.Equal(true, archived[0].Archived)
}

func TestFilterUnarchivedByDefault(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	list.Load(store.Load())
	filter := NewFilter(list.Todos())
	unarchived := filter.filterArchived("l")
	assert.Equal(1, len(unarchived))
	assert.Equal(false, unarchived[0].Archived)
}

func TestGetArchived(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	list.Load(store.Load())
	filter := NewFilter(list.Todos())
	archived := filter.getArchived()
	assert.Equal(1, len(archived))
	assert.Equal(true, archived[0].Archived)
}

func TestGetUnarchived(t *testing.T) {
	assert := assert.New(t)
	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	list.Load(store.Load())
	filter := NewFilter(list.Todos())
	unarchived := filter.getUnarchived()
	assert.Equal(1, len(unarchived))
	assert.Equal(false, unarchived[0].Archived)
}
