package ultralist

import "time"

// EditTodo edits a todo based upon a filter
func EditTodo(todo *Todo, filter *Filter) error {
	if filter.HasDue {
		dateParser := &DateParser{}
		dueDate, err := dateParser.ParseDate(filter.LastDue(), time.Now())
		if err != nil {
			return err
		}

		if dueDate.IsZero() {
			todo.Due = ""
		} else {
			todo.Due = dueDate.Format("2006-01-02")
		}
	}

	if filter.HasCompleted {
		if filter.Completed {
			todo.Complete()
		} else {
			todo.Uncomplete()
		}
	}

	if filter.HasArchived {
		todo.Archived = filter.Archived
	}

	if filter.HasIsPriority {
		todo.IsPriority = filter.IsPriority
	}

	if filter.HasStatus {
		todo.Status = filter.LastStatus()
	}

	if filter.Subject != "" {
		todo.Subject = filter.Subject
	}

	return nil
}
