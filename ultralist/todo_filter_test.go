package ultralist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoFilterCriteriaExcludesArchived(t *testing.T) {
	assert := assert.New(t)

	todoFilter := &TodoFilter{
		Filter: &Filter{},
		Todos:  SetupTodoList(),
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(5, len(filtered))
	for _, todo := range filtered {
		assert.Equal(false, todo.Archived)
	}
}

func TestFilterPriority(t *testing.T) {
	assert := assert.New(t)

	todoFilter := &TodoFilter{
		Filter: &Filter{HasIsPriority: true, IsPriority: true},
		Todos:  SetupTodoList(),
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(1, len(filtered))
	assert.Equal("has priority", filtered[0].Subject)

	todoFilter = &TodoFilter{
		Filter: &Filter{HasIsPriority: true, IsPriority: false},
		Todos:  SetupTodoList(),
	}

	filtered = todoFilter.ApplyFilter()

	assert.Equal(4, len(filtered))
	assert.Equal("not priority", filtered[0].Subject)
}

func TestFilterInclusive(t *testing.T) {
	assert := assert.New(t)

	todoFilter := &TodoFilter{
		Filter: &Filter{
			HasProjectFilter: true,
			Projects:         []string{"p3"},
		},
		Todos: SetupTodoList(),
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(2, len(filtered))
	assert.Equal("+p1 +p2 +p3", filtered[0].Subject)
	assert.Equal("+p1 +p3", filtered[1].Subject)
}

func TestFilterExclusive(t *testing.T) {
	assert := assert.New(t)

	todoFilter := &TodoFilter{
		Filter: &Filter{
			HasProjectFilter: true,
			ExcludeProjects:  []string{"p2"},
		},
		Todos: SetupTodoList(),
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(4, len(filtered))
	assert.Equal("has priority", filtered[0].Subject)
	assert.Equal("not priority", filtered[1].Subject)
	assert.Equal("+p1 +p3", filtered[2].Subject)
	assert.Equal("+p1", filtered[3].Subject)
}

func TestFilterInclusveAndExclusive(t *testing.T) {
	assert := assert.New(t)

	todoFilter := &TodoFilter{
		Filter: &Filter{
			HasProjectFilter: true,
			Projects:         []string{"p1"},
			ExcludeProjects:  []string{"p2"},
		},
		Todos: SetupTodoList(),
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(2, len(filtered))
	assert.Equal("+p1 +p3", filtered[0].Subject)
	assert.Equal("+p1", filtered[1].Subject)
}

func SetupTodoList() []*Todo {
	var todos []*Todo

	parser := &InputParser{}

	filter, _ := parser.Parse("has priority priority:true")
	todo, _ := CreateTodo(filter)
	todos = append(todos, todo)

	filter, _ = parser.Parse("not priority priority:false")
	todo, _ = CreateTodo(filter)
	todos = append(todos, todo)

	filter, _ = parser.Parse("archived archived:true")
	todo, _ = CreateTodo(filter)
	todos = append(todos, todo)

	filter, _ = parser.Parse("+p1 +p2 +p3")
	todo, _ = CreateTodo(filter)
	todos = append(todos, todo)

	filter, _ = parser.Parse("+p1 +p3")
	todo, _ = CreateTodo(filter)
	todos = append(todos, todo)

	filter, _ = parser.Parse("+p1")
	todo, _ = CreateTodo(filter)
	todos = append(todos, todo)

	return todos

}
