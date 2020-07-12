package ultralist

import "regexp"

// TodoFilter filters todos based on patterns.
type TodoFilter struct {
	Todos []*Todo
}

// NewFilter starts a todo filter and filters todos based on patterns.
func NewFilter(todos []*Todo) *TodoFilter {
	return &TodoFilter{Todos: todos}
}

// Filter filters todos based on patterns.
func (f *TodoFilter) Filter(input string) []*Todo {
	f.Todos = f.filterArchived(input)
	f.Todos = f.filterCompleted(input)
	f.Todos = f.filterPrioritized(input)
	f.Todos = f.filterStarted(input)
	f.Todos = f.filterProjects(input)
	f.Todos = f.filterContexts(input)
	f.Todos = NewDateFilter(f.Todos).FilterDate(input)

	return f.Todos
}

func (f *TodoFilter) isFilteringByProjects(input string) bool {
	parser := &Parser{}
	return len(parser.Projects(input)) > 0
}

func (f *TodoFilter) isFilteringByContexts(input string) bool {
	parser := &Parser{}
	return len(parser.Contexts(input)) > 0
}

func (f *TodoFilter) filterArchived(input string) []*Todo {

	r, _ := regexp.Compile(`is:archived`)
	if r.MatchString(input) {
		return f.getArchived()
	}

	return f.getUnarchived()
}

func (f *TodoFilter) filterCompleted(input string) []*Todo {

	completedRegex, _ := regexp.Compile(`is:completed`)
	if completedRegex.MatchString(input) {
		return f.getCompleted()
	}

	notCompletedRegex, _ := regexp.Compile(`not:completed`)
	if notCompletedRegex.MatchString(input) {
		return f.getNotCompleted()
	}

	return f.Todos
}

func (f *TodoFilter) filterPrioritized(input string) []*Todo {
	prioritizedRegex, _ := regexp.Compile(`is:priority`)
	if prioritizedRegex.MatchString(input) {
		return f.getPrioritized()
	}

	notPrioritizedRegex, _ := regexp.Compile(`not:priority`)
	if notPrioritizedRegex.MatchString(input) {
		return f.getNotPrioritized()
	}

	return f.Todos
}

func (f *TodoFilter) filterStarted(input string) []*Todo {
	prioritizedRegex, _ := regexp.Compile(`is:started`)
	if prioritizedRegex.MatchString(input) {
		return f.getStarted()
	}

	notPrioritizedRegex, _ := regexp.Compile(`not:started`)
	if notPrioritizedRegex.MatchString(input) {
		return f.getNotStarted()
	}

	return f.Todos
}

func (f *TodoFilter) filterProjects(input string) []*Todo {
	if !f.isFilteringByProjects(input) {
		return f.Todos
	}
	parser := &Parser{}
	projects := parser.Projects(input)
	var ret []*Todo

	for _, todo := range f.Todos {
		for _, todoProject := range todo.Projects {
			for _, project := range projects {
				if project == todoProject {
					ret = AddTodoIfNotThere(ret, todo)
				}
			}
		}
	}
	return ret
}

func (f *TodoFilter) filterContexts(input string) []*Todo {
	if !f.isFilteringByContexts(input) {
		return f.Todos
	}
	parser := &Parser{}
	contexts := parser.Contexts(input)
	var ret []*Todo

	for _, todo := range f.Todos {
		for _, todoContext := range todo.Contexts {
			for _, context := range contexts {
				if context == todoContext {
					ret = AddTodoIfNotThere(ret, todo)
				}
			}
		}
	}
	return ret
}

func (f *TodoFilter) getArchived() []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if todo.Archived {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *TodoFilter) getCompleted() []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if todo.Completed {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *TodoFilter) getNotCompleted() []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if !todo.Completed {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *TodoFilter) getPrioritized() []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if todo.IsPriority {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *TodoFilter) getNotPrioritized() []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if !todo.IsPriority {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *TodoFilter) getStarted() []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if todo.StartedDate != "" {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *TodoFilter) getNotStarted() []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if todo.StartedDate == "" {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *TodoFilter) getUnarchived() []*Todo {
	var ret []*Todo
	for _, todo := range f.Todos {
		if !todo.Archived {
			ret = append(ret, todo)
		}
	}
	return ret
}
