package todolist

type Todo struct {
	Id        int      `json:"id"`
	Subject   string   `json:"subject"`
	Projects  []string `json:"projects"`
	Contexts  []string `json:"contexts"`
	Due       string   `json:"due"`
	Completed bool     `json:"completed"`
	Archived  bool     `json:"archived"`
}

func NewTodo() *Todo {
	return &Todo{Completed: false, Archived: false}
}

func (t Todo) Valid() bool {
	return (t.Subject != "")
}
