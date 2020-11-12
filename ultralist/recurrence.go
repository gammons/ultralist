package ultralist

import (
	"time"
)

// Recurrence struct contains the logic for dealing with recurring todos.
type Recurrence struct{}

// HasNextRecurringTodo determines if a todo has a next recurrence.
func (r *Recurrence) HasNextRecurringTodo(todo *Todo) bool {
	recurUntil, _ := time.Parse(DATE_FORMAT, todo.RecurUntil)
	dueDate, _ := time.Parse(DATE_FORMAT, todo.Due)

	if todo.Recur != "" && todo.RecurUntil == "" {
		return true
	} else {
		return todo.Recur != "" && r.nextRecurrence(dueDate, time.Now(), todo.Recur).Before(recurUntil.AddDate(0, 0, 1))
	}
}

// NextRecurringTodo generates the next recurring todo from the one passed in.
func (r *Recurrence) NextRecurringTodo(todo *Todo, completedDate time.Time) *Todo {
	dueDate, _ := time.Parse(DATE_FORMAT, todo.Due)
	nextDueDate := r.nextRecurrence(dueDate, completedDate, todo.Recur)

	return &Todo{
		UUID:       newUUID(),
		Completed:  false,
		Archived:   false,
		Subject:    todo.Subject,
		Projects:   todo.Projects,
		Contexts:   todo.Contexts,
		Status:     todo.Status,
		IsPriority: todo.IsPriority,
		Notes:      todo.Notes,
		Recur:      todo.Recur,
		Due:        nextDueDate.Format(DATE_FORMAT),
		RecurUntil: todo.RecurUntil,
	}
}

func (r *Recurrence) nextRecurrence(dueDate time.Time, completedDate time.Time, recurrence string) time.Time {
	switch recurrence {
	case "daily":
		if completedDate.Before(dueDate) {
			return dueDate.AddDate(0, 0, 1)
		}
		return completedDate.AddDate(0, 0, 1)
	case "weekdays":
		return r.findNextWeekDay(dueDate, completedDate)
	case "weekly":
		return r.findNextWeek(dueDate, completedDate)
	case "monthly":
		return r.findNextMonth(dueDate, completedDate)
	case "yearly":
		return r.findNextYear(dueDate, completedDate)
	}
	return dueDate
}

func (r *Recurrence) findNextWeekDay(dueDate time.Time, completedDate time.Time) time.Time {
	dueDate = dueDate.AddDate(0, 0, 1)

	for {
		if !r.isWeekday(dueDate) || dueDate.Before(completedDate) {
			dueDate = dueDate.AddDate(0, 0, 1)
		} else {
			return dueDate
		}
	}
}

func (r *Recurrence) isWeekday(t time.Time) bool {
	switch t.Weekday() {
	case
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday:
		return true
	}
	return false
}

func (r *Recurrence) findNextWeek(dueDate time.Time, completedDate time.Time) time.Time {
	weekday := dueDate.Weekday()
	dueDate = dueDate.AddDate(0, 0, 1)
	for {
		if dueDate.Weekday() != weekday || dueDate.Before(completedDate) {
			dueDate = dueDate.AddDate(0, 0, 1)
		} else {
			return dueDate
		}
	}
}

func (r *Recurrence) findNextMonth(dueDate time.Time, completedDate time.Time) time.Time {
	dueDate = dueDate.AddDate(0, 1, 0)
	for {
		if dueDate.Before(completedDate) {
			dueDate = dueDate.AddDate(0, 1, 0)
		} else {
			return dueDate
		}
	}
}

func (r *Recurrence) findNextYear(dueDate time.Time, completedDate time.Time) time.Time {
	dueDate = dueDate.AddDate(1, 0, 0)
	for {
		if dueDate.Before(completedDate) {
			dueDate = dueDate.AddDate(1, 0, 0)
		} else {
			return dueDate
		}
	}
}
