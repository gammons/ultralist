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

func TestAddDoneTodo(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}

	app.AddDoneTodo("Groked how to do done todos @pop")

	todo := app.TodoList.FindById(1)
	assert.Equal("Groked how to do done todos @pop", todo.Subject)
	assert.Equal(true, todo.Completed)
	assert.Equal(false, todo.Archived)
	assert.Equal(false, todo.IsPriority)
	assert.Equal([]string{}, todo.Projects)
	assert.Equal(1, len(todo.Contexts))
	assert.Equal("pop", todo.Contexts[0])
}

func TestAddTodoWithEuropeanDates(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	year := strconv.Itoa(time.Now().Year())

	app.AddTodo("a do some stuff due 23 may")

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

func TestAddEmptyTodo(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}

	app.AddTodo("a")
	app.AddTodo("a      ")
	app.AddTodo("a\t\t\t\t")
	app.AddTodo("a\t \t  \t   \t")

	assert.Equal(len(app.TodoList.Data), 0)
}

func TestListbyProject(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	app.Load()

	// create three todos w/wo a project
	app.AddTodo("this is a test +testme")
	app.AddTodo("this is a test +testmetoo @work")
	app.AddTodo("this is a test with no projects")
	app.CompleteTodo("c 1")

	// simulate listTodos
	input := "l by p"
	filtered := NewFilter(app.TodoList.Todos()).Filter(input)
	grouped := app.getGroups(input, filtered)

	assert.Equal(3, len(grouped.Groups))

	// testme project has 1 todo and its completed
	assert.Equal(1, len(grouped.Groups["testme"]))
	assert.Equal(true, grouped.Groups["testme"][0].Completed)

	// testmetoo project has 1 todo and it has a context
	assert.Equal(1, len(grouped.Groups["testmetoo"]))
	assert.Equal(1, len(grouped.Groups["testmetoo"][0].Contexts))
	assert.Equal("work", grouped.Groups["testmetoo"][0].Contexts[0])
}

func TestListbyContext(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	app.Load()

	// create three todos w/wo a context
	app.AddTodo("this is a test +testme")
	app.AddTodo("this is a test +testmetoo @work")
	app.AddTodo("this is a test with no projects")
	app.CompleteTodo("c 1")

	// simulate listTodos
	input := "l by c"
	filtered := NewFilter(app.TodoList.Todos()).Filter(input)
	grouped := app.getGroups(input, filtered)

	assert.Equal(2, len(grouped.Groups))

	// work context has 1 todo and it has a project of testmetoo
	assert.Equal(1, len(grouped.Groups["work"]))
	assert.Equal(1, len(grouped.Groups["work"][0].Projects))
	assert.Equal("testmetoo", grouped.Groups["work"][0].Projects[0])

	// There are two todos with no context
	assert.Equal(2, len(grouped.Groups["No contexts"]))

	// check to see if the a todos with no context contain a
	// completed todo
	var hasACompletedTodo bool
	for _, todo := range grouped.Groups["No contexts"] {
		if todo.Completed {
			hasACompletedTodo = true
		}
	}
	assert.Equal(true, hasACompletedTodo)
}

func TestGetId(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	// not a valid id
	assert.Equal(-1, app.getId("p"))
	// a single digit id
	assert.Equal(6, app.getId("6"))
	// a double digit id
	assert.Equal(66, app.getId("66"))
}

func TestGetIds(t *testing.T) {
	assert := assert.New(t)
	app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	// no valid id here
	assert.Equal(0, len(app.getIds("p")))
	// one valid value here
	assert.Equal([]int{6}, app.getIds("6"))
	// lots of single post numbers
	assert.Equal([]int{6, 10, 8, 4}, app.getIds("6,10,8,4"))
	// a correct range
	assert.Equal([]int{6, 7, 8}, app.getIds("6-8"))
	// some incorrect ranges
	assert.Equal(0, len(app.getIds("6-6")))
	assert.Equal(0, len(app.getIds("8-6")))
	// some compsite ranges
	assert.Equal([]int{5, 6, 7, 8, 10, 11, 9}, app.getIds("5,6-8,10-11,9"))
}
