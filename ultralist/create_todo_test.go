package ultralist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateCompletedTodo(t *testing.T) {
	assert := assert.New(t)

	parser := &InputParser{}
	filter, _ := parser.Parse("lunch with bob completed:true")

	filter = &Filter{
		HasCompleted: true,
		Completed: true,
	}

	todo, _ := CreateTodo(filter)
	assert.Equal("lunch with bob", todo.Subject)
	assert.Equal("", todo.Due)
	assert.Equal(false, todo.IsPriority)
	assert.Equal(false, todo.Archived)
	assert.Equal(true, todo.Completed)
	assert.Equal([]string{}, todo.Projects)
	assert.Equal([]string{}, todo.Contexts)
	assert.Equal("completed", todo.Status)
}
