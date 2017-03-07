package todolist

import "time"

type Todo struct {
	Id         int      `json:"id"`
	Subject    string   `json:"subject"`
	Projects   []string `json:"projects"`
	Contexts   []string `json:"contexts"`
	Due        string   `json:"due"`
	Completed  bool     `json:"completed"`
	Archived   bool     `json:"archived"`
	IsPriority bool     `json:"isPriority"`
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
