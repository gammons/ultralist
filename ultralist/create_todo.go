package ultralist

// CreateTodo will create a TodoItem from a Filter
func CreateTodo(filter *Filter) (*Todo, error) {
	todoItem := &Todo{
		Subject:    filter.Subject,
		Archived:   filter.Archived,
		IsPriority: filter.IsPriority,
		Completed:  filter.Completed,
		Projects:   filter.Projects,
		Contexts:   filter.Contexts,
		Due:        filter.Due,
		Status:     filter.LastStatus(),
		Recur:      filter.Recur,
		RecurUntil: filter.RecurUntil,
	}
	if todoItem.Completed {
		todoItem.Complete()
	}

	return todoItem, nil
}
