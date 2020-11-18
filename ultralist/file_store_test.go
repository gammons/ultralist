package ultralist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStore(t *testing.T) {
	assert := assert.New(t)
	list := SetUpTestMemoryTodoList()
	store := &FileStore{}
	defer testFileCleanUp()
	list.FindByID(2).Subject = "this is an non-fixture subject"
	store.Initialize()
	store.Save(list.Todos())

	store1 := &FileStore{}

	todos, _ := store1.Load()
	assert.Equal(todos[1].Subject, "this is the first subject", "")
	assert.Equal(todos[0].Subject, "this is an non-fixture subject", "")
}
