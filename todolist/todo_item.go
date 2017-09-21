package todolist

import "time"

// Timestamp format to include date, time with timezone support. Easy to parse
const ISO8601_TIMESTAMP_FORMAT = "2006-01-02T15:04:05Z07:00"

type Todo struct {
	Id            int      `json:"id"`
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

func NewTodo() *Todo {
	return &Todo{Completed: false, Archived: false, IsPriority: false}
}

func (t Todo) Valid() bool {
	return (t.Subject != "")
}

func (t Todo) CalculateDueTime() time.Time {
	if t.Due != "" {
		parsedTime, _ := time.Parse("2006-01-02", t.Due)
		return parsedTime
	} else {
		parsedTime, _ := time.Parse("2006-01-02", "1900-01-01")
		return parsedTime
	}
}

func (t *Todo) Complete() {
	t.Completed = true
	t.CompletedDate = timestamp(time.Now()).Format(ISO8601_TIMESTAMP_FORMAT)
}

func (t *Todo) Uncomplete() {
	t.Completed = false
	t.CompletedDate = ""
}

func (t *Todo) Archive() {
	t.Archived = true
}

func (t *Todo) Unarchive() {
	t.Archived = false
}

func (t *Todo) Prioritize() {
	t.IsPriority = true
}

func (t *Todo) Unprioritize() {
	t.IsPriority = false
}

func (t Todo) CompletedDateToDate() string {
	parsedTime, _ := time.Parse(ISO8601_TIMESTAMP_FORMAT, t.CompletedDate)
	return parsedTime.Format("2006-01-02")
}
