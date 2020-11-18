package ultralist

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Integration-level tests for ultralist app as a whole

func TestAddTodo(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	year := strconv.Itoa(time.Now().Year())

	app.AddTodo("do some stuff due:dec31")

	todo := app.TodoList.FindByID(1)
	assert.Equal("do some stuff", todo.Subject)
	assert.Equal(fmt.Sprintf("%s-12-31", year), todo.Due)
	assert.Equal(false, todo.Completed)
	assert.Equal(false, todo.Archived)
	assert.Equal(false, todo.IsPriority)
	assert.Equal("", todo.CompletedDate)
	assert.Equal([]string{}, todo.Projects)
	assert.Equal([]string{}, todo.Contexts)
}

func TestAddTodoNoDue(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}

	app.AddTodo("do some stuff")

	todo := app.TodoList.FindByID(1)
	assert.Equal("do some stuff", todo.Subject)
	assert.Equal("", todo.Due)
	assert.Equal(false, todo.Completed)
	assert.Equal(false, todo.Archived)
	assert.Equal(false, todo.IsPriority)
	assert.Equal("", todo.CompletedDate)
	assert.Equal([]string{}, todo.Projects)
	assert.Equal([]string{}, todo.Contexts)
}

func TestGetId(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	// not a valid id
	id, err := app.getID("p")
	assert.Equal(-1, id)
	assert.Equal(true, err != nil)
	// a single digit id
	id, err = app.getID("6")
	assert.Equal(6, id)
	assert.Equal(true, err == nil)

	// a double digit id
	id, err = app.getID("66")
	assert.Equal(66, id)
	assert.Equal(true, err == nil)
}

func TestGetIds(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	// no valid id here
	assert.Equal(0, len(app.getIDs("p")))
	// one valid value here
	assert.Equal([]int{6}, app.getIDs("6"))
	// lots of single post numbers
	assert.Equal([]int{6, 10, 8, 4}, app.getIDs("6,10,8,4"))
	// a correct range
	assert.Equal([]int{6, 7, 8}, app.getIDs("6-8"))
	// some incorrect ranges
	assert.Equal(0, len(app.getIDs("6-6")))
	assert.Equal(0, len(app.getIDs("8-6")))
	// some compsite ranges
	assert.Equal([]int{5, 6, 7, 8, 10, 11, 9}, app.getIDs("5,6-8,10-11,9"))
}
