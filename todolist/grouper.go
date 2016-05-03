package todolist

type Grouper struct{}

type GroupedTodos struct {
	Groups map[string][]*Todo
}

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

func (g *Grouper) GroupByNothing(todos []*Todo) *GroupedTodos {
	groups := map[string][]*Todo{}
	groups["all"] = todos
	return &GroupedTodos{Groups: groups}
}
