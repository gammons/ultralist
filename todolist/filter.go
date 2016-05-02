package todolist

import "regexp"

type TodoFilter struct {
	Todos []Todo
}

func NewFilter(todos []Todo) *TodoFilter {
	return &TodoFilter{Todos: todos}
}

func (f *TodoFilter) Filter(input string) []Todo {
	f.Todos = f.filterArchived(input)
	return f.Todos
}

func (f *TodoFilter) filterArchived(input string) []Todo {
	r, _ := regexp.Compile(`l archived$`)
	if r.MatchString(input) {
		return f.getArchived()
	} else {
		return f.getUnarchived()
	}
}

func (f *TodoFilter) getArchived() []Todo {
	var ret []Todo
	for _, todo := range f.Todos {
		if todo.Archived == true {
			ret = append(ret, todo)
		}
	}
	return ret
}

func (f *TodoFilter) getUnarchived() []Todo {
	var ret []Todo
	for _, todo := range f.Todos {
		if todo.Archived == false {
			ret = append(ret, todo)
		}
	}
	return ret
}
