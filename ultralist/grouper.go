package ultralist

// Grouper is the group struct.
type Grouper struct{}

// GroupedTodos is the main struct of this file.
type GroupedTodos struct {
	Groups map[string][]*Todo
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
	return &GroupedTodos{Groups: groups}
}

// GroupByNothing is the default result if todos are not grouped by context project.
func (g *Grouper) GroupByNothing(todos []*Todo) *GroupedTodos {
	groups := map[string][]*Todo{}
	groups["all"] = todos
	return &GroupedTodos{Groups: groups}
}
