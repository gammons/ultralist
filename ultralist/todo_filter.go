package ultralist

import (
	"time"
)

// TodoFilter filters todos based on patterns.
type TodoFilter struct {
	Filter *Filter
	Todos  []*Todo
}

// ApplyFilter filters todos based on the Filter struct passed in.
func (f *TodoFilter) ApplyFilter() []*Todo {
	var filtered []*Todo

	for _, todo := range f.Todos {
		if f.Filter.HasIsPriority {
			if todo.IsPriority != f.Filter.IsPriority {
				continue
			}
		}

		if f.Filter.HasCompleted {
			if todo.Completed != f.Filter.Completed {
				continue
			}
		}

		if f.Filter.HasArchived {
			if todo.Archived != f.Filter.Archived {
				continue
			}
		}

		if f.Filter.HasStatus {
			if !f.todoPassesFilter([]string{todo.Status}, f.Filter.Status, f.Filter.ExcludeStatus) {
				continue
			}
		}

		if f.Filter.HasProjectFilter {
			if !f.todoPassesFilter(todo.Projects, f.Filter.Projects, f.Filter.ExcludeProjects) {
				continue
			}
		}

		if f.Filter.HasContextFilter {
			if !f.todoPassesFilter(todo.Contexts, f.Filter.Contexts, f.Filter.ExcludeContexts) {
				continue
			}
		}

		// has exact due date
		if f.Filter.HasDue {
			if todo.Due != f.Filter.Due {
				continue
			}
		}

		if f.Filter.HasDueBefore {
			todoTime, _ := time.Parse(DATE_FORMAT, todo.Due)
			dueBeforeTime, _ := time.Parse(DATE_FORMAT, f.Filter.DueBefore)

			if todo.Due == "" || !todoTime.Before(dueBeforeTime) {
				continue
			}
		}

		if f.Filter.HasDueAfter {
			todoTime, _ := time.Parse(DATE_FORMAT, todo.Due)
			dueAfterTime, _ := time.Parse(DATE_FORMAT, f.Filter.DueAfter)

			if todo.Due == "" || !todoTime.After(dueAfterTime) {
				continue
			}
		}

		if f.Filter.HasRecur {
			if todo.Recur != f.Filter.Recur {
				continue
			}
		}

		// finally, if we're still here, append the todo to the filtered list
		filtered = append(filtered, todo)
	}

	// the "default" filter is to filter out archived todos, if nothing is set.
	if !f.Filter.HasArchived {
		var ret []*Todo

		for _, todo := range filtered {
			if !todo.Archived {
				ret = append(ret, todo)
			}
		}
		return ret
	}

	return filtered
}

func (f *TodoFilter) todoPassesFilter(val []string, inclusiveVals []string, exclusiveVals []string) bool {
	ret := true
	// inclusive vals is evaluated via OR
	if len(inclusiveVals) >= 1 {
		ret = false

		for _, iv := range inclusiveVals {
			for _, v := range val {
				if iv == v {
					ret = true
				}
			}
		}
	}

	// exclusiveVals is evaluated with AND
	for _, ev := range exclusiveVals {
		for _, v := range val {
			if ev == v {
				ret = false
			}
		}
	}

	return ret
}
