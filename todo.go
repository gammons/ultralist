package todolist

import "time"

type Todo struct {
	Id        string
	Subject   string
	Projects  []string
	Contexts  []string
	Due       time.Time
	Completed bool
	Archived  bool
}

func NewTodo() *Todo {
	return &Todo{Completed: false, Archived: false}
}

func (t Todo) Valid() bool {
	return (t.Subject != "")
}
