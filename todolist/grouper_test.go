package todolist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupByContext(t *testing.T) {
	assert := assert.New(t)

	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	todos, _ := store.Load()
	list.Load(todos)

	grouper := &Grouper{}
	grouped := grouper.GroupByContext(list.Todos())

	assert.Equal(2, len(grouped.Groups["root"]), "")
	assert.Equal(1, len(grouped.Groups["more"]), "")
}

func TestGroupByProject(t *testing.T) {
	assert := assert.New(t)

	store := &FileStore{FileLocation: "todos.json"}
	list := &TodoList{}
	todos, _ := store.Load()
	list.Load(todos)

	grouper := &Grouper{}
	grouped := grouper.GroupByProject(list.Todos())

	assert.Equal(2, len(grouped.Groups["test1"]), "")
}
