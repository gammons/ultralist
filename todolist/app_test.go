package todolist

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddTodo(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	year := strconv.Itoa(time.Now().Year())

	app.AddTodo("a do some stuff due may 23")

	todo := app.TodoList.FindById(1)
	assert.Equal("do some stuff", todo.Subject)
	assert.Equal(fmt.Sprintf("%s-05-23", year), todo.Due)
	assert.Equal(false, todo.Completed)
	assert.Equal(false, todo.Archived)
	assert.Equal(false, todo.IsPriority)
	assert.Equal("", todo.CompletedDate)
	assert.Equal([]string{}, todo.Projects)
	assert.Equal([]string{}, todo.Contexts)
}

func TestAddTodoWithEuropeanDates(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}

	app.AddTodo("a do some stuff due 23 may")

	todo := app.TodoList.FindById(1)
	assert.Equal("do some stuff", todo.Subject)
	assert.Equal("2017-05-23", todo.Due)
	assert.Equal(false, todo.Completed)
	assert.Equal(false, todo.Archived)
	assert.Equal(false, todo.IsPriority)
	assert.Equal("", todo.CompletedDate)
	assert.Equal([]string{}, todo.Projects)
	assert.Equal([]string{}, todo.Contexts)
}
