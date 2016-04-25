package todolist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupByContext(t *testing.T) {
	assert := assert.New(t)

	store := &FileStore{FileLocation: "todos.json"}
	store.Load()

	grouper := &Grouper{}
	grouped := grouper.GroupByContext(store.Todos())

	assert.Equal(2, len(grouped.Groups["root"]), "")
	assert.Equal(1, len(grouped.Groups["more"]), "")
}

func TestGroupByProject(t *testing.T) {
	assert := assert.New(t)

	store := &FileStore{FileLocation: "todos.json"}
	store.Load()

	grouper := &Grouper{}
	grouped := grouper.GroupByProject(store.Todos())

	assert.Equal(2, len(grouped.Groups["test1"]), "")
}
