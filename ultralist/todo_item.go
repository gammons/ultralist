package ultralist

import (
	"reflect"
	"time"
)

// iso8601TimestampFormat is the timestamp format to include date, time with timezone support. Easy to parse.
const iso8601TimestampFormat = "2006-01-02T15:04:05Z07:00"

// Todo is the struct for a todo item.
type Todo struct {
	ID                int      `json:"id"`
	UUID              string   `json:"uuid"`
	Subject           string   `json:"subject"`
	Projects          []string `json:"projects"`
	Contexts          []string `json:"contexts"`
	Due               string   `json:"due"`
	Completed         bool     `json:"completed"`
	CompletedDate     string   `json:"completed_date"`
	Status            string   `json:"status"`
	Archived          bool     `json:"archived"`
	IsPriority        bool     `json:"is_priority"`
	Notes             []string `json:"notes"`
	Recur             string   `json:"recur"`
	RecurUntil        string   `json:"recur_until"`
	PrevRecurTodoUUID string   `json:"prev_recur_todo_uuid"`
}

// NewTodo is creating a new todo item.
func NewTodo() *Todo {
	return &Todo{UUID: newUUID(), Completed: false, Archived: false, IsPriority: false}
}

// Valid is checking if a new todo is valid or not.
func (t Todo) Valid() bool {
	return (t.Subject != "")
}

// CalculateDueTime is calculating the due time of the todo item.
func (t Todo) CalculateDueTime() time.Time {
	if t.Due != "" {
		parsedTime, _ := time.Parse(DATE_FORMAT, t.Due)
		return parsedTime
	}
	parsedTime, _ := time.Parse(DATE_FORMAT, "1900-01-01")
	return parsedTime
}

// Complete is completing a todo item and sets the complete date to the current time.
func (t *Todo) Complete() {
	t.Completed = true
	t.Status = "completed"
	t.CompletedDate = timestamp(time.Now()).Format(iso8601TimestampFormat)
}

// Uncomplete is uncompleting a todo item and removes the complete date.
func (t *Todo) Uncomplete() {
	t.Completed = false
	t.Status = ""
	t.CompletedDate = ""
}

// Archive is archiving a todo item.
func (t *Todo) Archive() {
	t.Archived = true
}

// Unarchive is unarchiving a todo item.
func (t *Todo) Unarchive() {
	t.Archived = false
}

// Prioritize is prioritizing a todo item.
func (t *Todo) Prioritize() {
	t.IsPriority = true
}

// Unprioritize is unpriotizing a todo item.
func (t *Todo) Unprioritize() {
	t.IsPriority = false
}

// CompletedDateToDate is returning the date when an item was completed.
func (t Todo) CompletedDateToDate() string {
	parsedTime, _ := time.Parse(iso8601TimestampFormat, t.CompletedDate)
	return parsedTime.Format(DATE_FORMAT)
}

// HasNotes is showing if an todo has notes.
func (t Todo) HasNotes() bool {
	return len(t.Notes) > 0
}

// Equals compares 2 todos for equality.
func (t Todo) Equals(other *Todo) bool {
	if t.ID != other.ID ||
		t.UUID != other.UUID ||
		t.Subject != other.Subject ||
		!reflect.DeepEqual(t.Projects, other.Projects) ||
		!reflect.DeepEqual(t.Contexts, other.Contexts) ||
		t.Due != other.Due ||
		t.Completed != other.Completed ||
		t.Status != other.Status ||
		t.CompletedDate != other.CompletedDate ||
		t.Archived != other.Archived ||
		t.IsPriority != other.IsPriority ||
		t.Recur != other.Recur ||
		t.RecurUntil != other.RecurUntil ||
		!reflect.DeepEqual(t.Notes, other.Notes) {
		return false
	}
	return true
}
