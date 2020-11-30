package ultralist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNoFilterCriteriaExcludesArchived(t *testing.T) {
	assert := assert.New(t)

	todos := SetupTodoList()
	todoFilter := &TodoFilter{
		Filter: &Filter{},
		Todos:  todos,
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(len(todos)-1, len(filtered))
	for _, todo := range filtered {
		assert.Equal(false, todo.Archived)
	}
}

func TestFilterPriority(t *testing.T) {
	assert := assert.New(t)

	todos := SetupTodoList()
	todoFilter := &TodoFilter{
		Filter: &Filter{HasIsPriority: true, IsPriority: true},
		Todos:  todos,
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(1, len(filtered))
	assert.Equal("has priority", filtered[0].Subject)

	todoFilter = &TodoFilter{
		Filter: &Filter{HasIsPriority: true, IsPriority: false},
		Todos:  SetupTodoList(),
	}

	filtered = todoFilter.ApplyFilter()

	assert.Equal(len(todos)-2, len(filtered))
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

	todos := SetupTodoList()
	todoFilter := &TodoFilter{
		Filter: &Filter{
			HasProjectFilter: true,
			ExcludeProjects:  []string{"p2"},
		},
		Todos: todos,
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(len(todos)-2, len(filtered))
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

func TestFilterDue(t *testing.T) {
	assert := assert.New(t)

	todoFilter := &TodoFilter{
		Filter: &Filter{
			HasDue: true,
			Due:    time.Now().Format(DateFormat),
		},
		Todos: SetupTodoList(),
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(1, len(filtered))
	assert.Equal("due today", filtered[0].Subject)
}

func TestFilterDueBefore(t *testing.T) {
	assert := assert.New(t)

	todoFilter := &TodoFilter{
		Filter: &Filter{
			HasDueBefore: true,
			DueBefore:    time.Now().Format(DateFormat),
		},
		Todos: SetupTodoList(),
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(1, len(filtered))
	assert.Equal("due yesterday", filtered[0].Subject)
}

func TestFilterDueAfter(t *testing.T) {
	assert := assert.New(t)

	todoFilter := &TodoFilter{
		Filter: &Filter{
			HasDueAfter: true,
			DueAfter:    time.Now().Format(DateFormat),
		},
		Todos: SetupTodoList(),
	}

	filtered := todoFilter.ApplyFilter()

	assert.Equal(1, len(filtered))
	assert.Equal("due tomorrow", filtered[0].Subject)
}

func SetupTodoList() []*Todo {
	var todos []*Todo

	todos = append(todos, &Todo{
		Subject:    "has priority",
		IsPriority: true,
	})

	todos = append(todos, &Todo{
		Subject:    "not priority",
		IsPriority: false,
	})

	todos = append(todos, &Todo{
		Subject:  "archived",
		Archived: true,
	})

	todos = append(todos, &Todo{
		Subject:  "+p1 +p2 +p3",
		Projects: []string{"p1", "p2", "p3"},
	})

	todos = append(todos, &Todo{
		Subject:  "+p1 +p3",
		Projects: []string{"p1", "p3"},
	})

	todos = append(todos, &Todo{
		Subject:  "+p1",
		Projects: []string{"p1"},
	})

	todos = append(todos, &Todo{
		Subject: "due today",
		Due:     time.Now().Format(DateFormat),
	})

	todos = append(todos, &Todo{
		Subject: "due tomorrow",
		Due:     time.Now().AddDate(0, 0, 1).Format(DateFormat),
	})

	todos = append(todos, &Todo{
		Subject: "due yesterday",
		Due:     time.Now().AddDate(0, 0, -1).Format(DateFormat),
	})

	return todos
}
