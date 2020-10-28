package ultralist

// EditTodo edits a todo based upon a filter
func EditTodo(todo *Todo, filter *Filter) error {
	if filter.HasDue {
		todo.Due = filter.Due
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

	if len(filter.Projects) > 0 {
		todo.Projects = filter.Projects
	}

	if len(filter.Contexts) > 0 {
		todo.Contexts = filter.Contexts
	}

	return nil
}
