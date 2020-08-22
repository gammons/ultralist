package ultralist

import "time"

// CreateTodo will create a TodoItem from a Filter
func CreateTodo(filter *Filter) (*Todo, error) {
	dueDateString := ""

	if filter.LastDue() != "" {
		dateParser := &DateParser{}
		dueDate, err := dateParser.ParseDate(filter.LastDue(), time.Now())
		if err != nil {
			return nil, err
		}
		dueDateString = dueDate.Format("2006-01-02")
	}

	todoItem := &Todo{
		Subject:    filter.Subject,
		Archived:   filter.Archived,
		IsPriority: filter.IsPriority,
		Completed:  filter.Completed,
		Projects:   filter.Projects,
		Contexts:   filter.Contexts,
		Due:        dueDateString,
		Status:     filter.LastStatus(),
	}
	if todoItem.Completed {
		todoItem.CompletedDate = time.Now().Format("2006-01-02")
		todoItem.Status = "completed"
	}

	return todoItem, nil
}
