package ultralist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditTodoProjects(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("+p1")
	todo, _ := CreateTodo(filter)
	todoList := &TodoList{}
	todoList.Data = []*Todo{todo}

	editFilter, _ := parser.Parse("+p2")

	EditTodo(todo, todoList, editFilter)

	assert.Equal("+p2", todo.Subject)
	assert.Equal([]string{"p2"}, todo.Projects)
}

func TestEditTodoProjectsOtherSyntax(t *testing.T) {
	assert := assert.New(t)
	parser := &InputParser{}

	filter, _ := parser.Parse("+p1")
	todo, _ := CreateTodo(filter)
	todoList := &TodoList{}
	todoList.Data = []*Todo{todo}

	editFilter, _ := parser.Parse("project:p2")

	EditTodo(todo, todoList, editFilter)

	assert.Equal("+p1", todo.Subject)
	assert.Equal([]string{"p2"}, todo.Projects)
}
