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

func TestAddTodoWithEuropeanDates(t *testing.T) {
	// assert := assert.New(t)
	// app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	// year := strconv.Itoa(time.Now().Year())
	//
	// app.AddTodo("a do some stuff due 31 dec")
	//
	// todo := app.TodoList.FindByID(1)
	// assert.Equal("do some stuff", todo.Subject)
	// assert.Equal(fmt.Sprintf("%s-12-31", year), todo.Due)
	// assert.Equal(false, todo.Completed)
	// assert.Equal(false, todo.Archived)
	// assert.Equal(false, todo.IsPriority)
	// assert.Equal("", todo.CompletedDate)
	// assert.Equal([]string{}, todo.Projects)
	// assert.Equal([]string{}, todo.Contexts)
}

func TestListbyProject(t *testing.T) {
	// assert := assert.New(t)
	// app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	// app.Load()
	//
	// // create three todos w/wo a project
	// app.AddTodo("this is a test +testme")
	// app.AddTodo("this is a test +testmetoo @work")
	// app.AddTodo("this is a test with no projects")
	// app.AddTodo("this is a test to archive completed todos automatically +testarchivecompleted")
	// app.CompleteTodo("c 1", false)
	// app.CompleteTodo("c 4", true)
	//
	// // simulate listTodos
	// input := "l group:p"
	// filtered := NewFilter(app.TodoList.Todos()).Filter(input)
	// grouped := app.getGroups(input, filtered)
	//
	// assert.Equal(3, len(grouped.Groups))
	//
	// // testme project has 1 todo and its completed
	// assert.Equal(1, len(grouped.Groups["testme"]))
	// assert.Equal(true, grouped.Groups["testme"][0].Completed)
	//
	// // testmetoo project has 1 todo and it has a context
	// assert.Equal(1, len(grouped.Groups["testmetoo"]))
	// assert.Equal(1, len(grouped.Groups["testmetoo"][0].Contexts))
	// assert.Equal("work", grouped.Groups["testmetoo"][0].Contexts[0])
	//
	// // testmetoo project has 0 todos
	// assert.Equal(0, len(grouped.Groups["testarchivecompleted"]))
}

func TestListbyContext(t *testing.T) {
	// assert := assert.New(t)
	// app := &App{TodoList: &TodoList{}, TodoStore: &MemoryStore{}}
	// app.Load()
	//
	// // create three todos w/wo a context
	// app.AddTodo("this is a test +testme")
	// app.AddTodo("this is a test +testmetoo @work")
	// app.AddTodo("this is a test with no projects")
	// app.CompleteTodo("c 1", false)
	//
	// // simulate listTodos
	// input := "l group:c"
	// filtered := NewFilter(app.TodoList.Todos()).Filter(input)
	// grouped := app.getGroups(input, filtered)
	//
	// assert.Equal(2, len(grouped.Groups))
	//
	// // work context has 1 todo and it has a project of testmetoo
	// assert.Equal(1, len(grouped.Groups["work"]))
	// assert.Equal(1, len(grouped.Groups["work"][0].Projects))
	// assert.Equal("testmetoo", grouped.Groups["work"][0].Projects[0])
	//
	// // There are two todos with no context
	// assert.Equal(2, len(grouped.Groups["No contexts"]))
	//
	// // check to see if the a todos with no context contain a
	// // completed todo
	// var hasACompletedTodo bool
	// for _, todo := range grouped.Groups["No contexts"] {
	// 	if todo.Completed {
	// 		hasACompletedTodo = true
	// 	}
	// }
	// assert.Equal(true, hasACompletedTodo)
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
