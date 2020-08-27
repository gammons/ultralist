package ultralist

import "time"

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
	}
	if todoItem.Completed {
		todoItem.CompletedDate = time.Now().Format("2006-01-02")
		todoItem.Status = "completed"
	}

	return todoItem, nil
}
