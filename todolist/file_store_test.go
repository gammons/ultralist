package todolist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStore(t *testing.T) {
	assert := assert.New(t)
	// load the fixutre data
	list := SetUpTestMemoryTodoList()
	testFilename := "TestFileStore_todos.json"
	store := &FileStore{FileLocation: testFilename}
	defer testFileCleanUp(testFilename)
	list.FindById(2).Subject = "this is an non-fixture subject"
	// save the fixture data
	store.Save(list.Todos())

	// reload the fixture data
	store1 := &FileStore{FileLocation: testFilename}
	todos, _ := store1.Load()
	assert.Equal(todos[1].Subject, "this is the first subject", "")
	assert.Equal(todos[0].Subject, "this is an non-fixture subject", "")
}
