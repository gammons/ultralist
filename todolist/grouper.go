package todolist

type Grouper struct{}

type GroupedTodos struct {
	Groups map[string][]Todo
}

func (g *Grouper) GroupByContext(todos []Todo) *GroupedTodos {
	groups := map[string][]Todo{}

	allContexts := []string{}

	for _, todo := range todos {
		allContexts = addIfNotThere(allContexts, todo.Contexts)
	}

	for _, todo := range todos {
		for _, context := range todo.Contexts {
			groups[context] = append(groups[context], todo)
		}
	}

	return &GroupedTodos{Groups: groups}
}

func (g *Grouper) GroupByProject(todos []Todo) *GroupedTodos {
	groups := map[string][]Todo{}

	allProjects := []string{}

	for _, todo := range todos {
		allProjects = addIfNotThere(allProjects, todo.Projects)
	}

	for _, todo := range todos {
		for _, project := range todo.Projects {
			groups[project] = append(groups[project], todo)
		}
	}
	return &GroupedTodos{Groups: groups}
}

func (g *Grouper) GroupByNothing(todos []Todo) *GroupedTodos {
	groups := map[string][]Todo{}
	groups["all"] = todos
	return &GroupedTodos{Groups: groups}
}

func addIfNotThere(arr []string, items []string) []string {
	for _, item := range items {
		there := false
		for _, arrItem := range arr {
			if item == arrItem {
				there = true
			}
		}
		if !there {
			arr = append(arr, item)
		}
	}
	return arr
}
