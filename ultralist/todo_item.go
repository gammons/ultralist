package ultralist

import (
	"reflect"
	"time"
)

// iso8601TimestampFormat is the timestamp format to include date, time with timezone support. Easy to parse.
const iso8601TimestampFormat = "2006-01-02T15:04:05Z07:00"

// Todo is the struct for a todo item.
type Todo struct {
	ID            int      `json:"id"`
	UUID          string   `json:"uuid"`
	Subject       string   `json:"subject"`
	Projects      []string `json:"projects"`
	Contexts      []string `json:"contexts"`
	Due           string   `json:"due"`
	Completed     bool     `json:"completed"`
	CompletedDate string   `json:"completedDate"`
	Archived      bool     `json:"archived"`
	IsPriority    bool     `json:"isPriority"`
	Notes         []string `json:"notes"`
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
		parsedTime, _ := time.Parse("2006-01-02", t.Due)
		return parsedTime
	}
	parsedTime, _ := time.Parse("2006-01-02", "1900-01-01")
	return parsedTime
}

// Complete is completing a todo item and sets the complete date to the current time.
func (t *Todo) Complete() {
	t.Completed = true
	t.CompletedDate = timestamp(time.Now()).Format(iso8601TimestampFormat)
}

// Uncomplete is uncompleting a todo item and removes the complete date.
func (t *Todo) Uncomplete() {
	t.Completed = false
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
	return parsedTime.Format("2006-01-02")
}

// set a Todo's subject and also set its projects and contexts.
func (t *Todo) SetSubject(subject string) {
	t.Subject = subject
	parser := &Parser{}
	t.Projects = parser.Projects(subject)
	t.Contexts = parser.Contexts(subject)
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
		t.CompletedDate != other.CompletedDate ||
		t.Archived != other.Archived ||
		t.IsPriority != other.IsPriority ||
		!reflect.DeepEqual(t.Notes, other.Notes) {
		return false
	}
	return true
}
