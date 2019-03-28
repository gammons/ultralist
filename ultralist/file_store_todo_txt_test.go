package ultralist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {

}

// func TestFileStoreTodoTxtLoad(t *testing.T) {
// 	// assert := assert.New(t)
//
// 	store := &FileStoreTodoTxt{}
// 	store.Load()
// }

func TestFileStoreTodoTxtParseLinePriority(t *testing.T) {
	assert := assert.New(t)
	store := &FileStoreTodoTxt{}
	todo := store.ParseLine("x (A) this is a priority completed item")
	assert.Equal(true, todo.IsPriority)

	todo = store.ParseLine("x this is not a prioirity")
	assert.Equal(false, todo.IsPriority)
}

func TestFileStoreTodoTxtParseLineCompleted(t *testing.T) {
	assert := assert.New(t)
	store := &FileStoreTodoTxt{}
	todo := store.ParseLine("x (A) this is a completed todo item")
	assert.Equal(true, todo.Completed)

	todo = store.ParseLine("x this is completed too")
	assert.Equal(true, todo.Completed)

	todo = store.ParseLine("(A) this is not completed")
	assert.Equal(false, todo.Completed)

	todo = store.ParseLine("this is not completed")
	assert.Equal(false, todo.Completed)
}

func TestFileStoreTodoTxtParseLineCompletedDate(t *testing.T) {
	assert := assert.New(t)
	store := &FileStoreTodoTxt{}
	todo := store.ParseLine("x 2018-05-01 this is a completed todo item")
	assert.Equal(true, todo.Completed)
	assert.Equal("2018-05-01", todo.CompletedDate)
}

func TestFileStoreTodoTxtParseLineDueDate(t *testing.T) {
	assert := assert.New(t)
	store := &FileStoreTodoTxt{}
	todo := store.ParseLine("x 2018-05-01 this is a completed todo item due:2018-06-01")
	assert.Equal("2018-06-01", todo.Due)

	todo = store.ParseLine("this has no due date")
	assert.Equal("", todo.Due)
}

func TestFileStoreTodoTxtParseLineID(t *testing.T) {
	assert := assert.New(t)
	store := &FileStoreTodoTxt{}
	todo := store.ParseLine("x 2018-05-01 this is a completed todo item due:2018-06-01 id:25")
	assert.Equal(25, todo.ID)
}

func TestFileStoreTodoTxtParseLineSubject(t *testing.T) {
	assert := assert.New(t)
	store := &FileStoreTodoTxt{}
	todo := store.ParseLine("x 2018-05-01 this is a completed todo item due:2018-06-01 id:25")
	assert.Equal("this is a completed todo item", todo.Subject)

	todo = store.ParseLine("here is a simple one")
	assert.Equal("here is a simple one", todo.Subject)

	todo = store.ParseLine("(A) prioritized todo")
	assert.Equal("prioritized todo", todo.Subject)
}
