package ultralist

import (
	"sort"
)

// Grouping dictates how todos should be grouped
type Grouping string

const (
	// ByNone dictates that todos should not be grouped at all.
	ByNone Grouping = "none"

	// ByContext dictates that todos should be grouped by their context.
	ByContext Grouping = "context"

	// ByProject dictates that todos should be grouped by their project.
	ByProject Grouping = "project"

	// ByStatus dictates that todos should be grouped by their status.
	ByStatus Grouping = "status"
)

// Grouper is the group struct.
type Grouper struct{}

// GroupedTodos is the main struct of this file.
type GroupedTodos struct {
	Groups map[string][]*Todo
}

// GroupTodos will group an array of todos by the specified Grouping.
func (g *Grouper) GroupTodos(todos []*Todo, grouping Grouping) *GroupedTodos {
	switch grouping {
	case ByContext:
		return g.GroupByContext(todos)
	case ByProject:
		return g.GroupByProject(todos)
	case ByStatus:
		return g.GroupByStatus(todos)
	}
	return g.GroupByNothing(todos)
}

// GroupByContext is grouping todos by its context.
func (g *Grouper) GroupByContext(todos []*Todo) *GroupedTodos {
	groups := map[string][]*Todo{}

	allContexts := []string{}

	for _, todo := range todos {
		allContexts = AddIfNotThere(allContexts, todo.Contexts)
	}

	for _, todo := range todos {
		for _, context := range todo.Contexts {
			groups[context] = append(groups[context], todo)
		}
		if len(todo.Contexts) == 0 {
			groups["No contexts"] = append(groups["No contexts"], todo)
		}
	}

	// finally, sort the todos
	for groupName, todos := range groups {
		groups[groupName] = g.sort(todos)
	}

	return &GroupedTodos{Groups: groups}
}

// GroupByProject is grouping todos by its project.
func (g *Grouper) GroupByProject(todos []*Todo) *GroupedTodos {
	groups := map[string][]*Todo{}

	allProjects := []string{}

	for _, todo := range todos {
		allProjects = AddIfNotThere(allProjects, todo.Projects)
	}

	for _, todo := range todos {
		for _, project := range todo.Projects {
			groups[project] = append(groups[project], todo)
		}
		if len(todo.Projects) == 0 {
			groups["No projects"] = append(groups["No projects"], todo)
		}
	}

	// finally, sort the todos
	for groupName, todos := range groups {
		groups[groupName] = g.sort(todos)
	}

	return &GroupedTodos{Groups: groups}
}

// GroupByStatus is grouping todos by status
func (g *Grouper) GroupByStatus(todos []*Todo) *GroupedTodos {
	groups := map[string][]*Todo{}

	for _, todo := range todos {
		if todo.Status != "" {
			groups[todo.Status] = append(groups[todo.Status], todo)
		} else {
			groups["No status"] = append(groups["No status"], todo)
		}
	}

	// finally, sort the todos
	for groupName, todos := range groups {
		groups[groupName] = g.sort(todos)
	}

	return &GroupedTodos{Groups: groups}
}

// GroupByNothing is the default result if todos are not grouped by context project.
func (g *Grouper) GroupByNothing(todos []*Todo) *GroupedTodos {
	groups := map[string][]*Todo{}
	groups["all"] = g.sort(todos)

	return &GroupedTodos{Groups: groups}
}

func (g *Grouper) sort(todos []*Todo) []*Todo {
	sort.Slice(todos, func(i, j int) bool {
		// always favor unarchived todos over archived ones
		if todos[i].Archived || todos[j].Archived {
			return !todos[i].Archived || todos[j].Archived
		}

		// always favor un-completed todos over completed ones
		if todos[i].Completed || todos[j].Completed {
			return !todos[i].Completed || todos[j].Completed
		}

		// always favor prioritized todos
		if todos[i].IsPriority || todos[j].IsPriority {
			return todos[i].IsPriority || !todos[j].IsPriority
		}

		// always prefer a todo with a due date
		if todos[i].Due != "" && todos[j].Due == "" {
			return true
		}

		// un-favor todos without a due date
		if todos[i].Due == "" {
			return false
		}

		return todos[i].CalculateDueTime().Before(todos[j].CalculateDueTime())
	})

	return todos
}
