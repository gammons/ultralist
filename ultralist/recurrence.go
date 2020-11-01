package ultralist

import (
	"time"
)

type Recurrence struct{}

// HasNextRecurringTodo
func (r *Recurrence) HasNextRecurringTodo(todo *Todo) bool {
	recurUntil, _ := time.Parse(DATE_FORMAT, todo.RecurUntil)

	return todo.Recur != "" &&
		time.Now().Before(recurUntil) &&
		!todo.Completed
}

// NextRecurringTodo generates the next recurring todo from the one passed in.
func (r *Recurrence) NextRecurringTodo(todo *Todo, completedTime time.Time) *Todo {

	dueDate, _ := time.Parse(DATE_FORMAT, todo.Due)
	var nextDueDate time.Time

	switch todo.Recur {
	case "daily":
		if completedTime.Before(dueDate) {
			nextDueDate = dueDate.AddDate(0, 0, 1)
		} else {
			nextDueDate = completedTime.AddDate(0, 0, 1)
		}
	case "weekdays":
		// m, t, w, thu, fri
		// schedule it for tomorrow unless it's currently friday or saturday, in which case schedule it for next monday
		nextDueDate = r.findNextWeekDay(dueDate, completedTime)
	case "monthly":
		// add one month to the due date
		// schedule the next on the next month day
		nextDueDate = r.findNextMonth(dueDate, completedTime)
	case "weekly":
		// add one week to the due date
		// schedule it for the next weekday that's the same day as the due date.
		// if the due date is in the past, schedule it for the next week day closest to today
		// if the due date is in the future, schedule it for the weekday after the due date

		nextDueDate = r.findNextWeek(dueDate, completedTime)
	}

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
