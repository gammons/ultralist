package todolist

type Store interface {
	Load()
	Save()
	Todos() []Todo

	Add(t *Todo)
	Delete(id int)

	Complete(id int)
	Uncomplete(id int)

	IndexOf(t *Todo) int
	FindById(id int) *Todo
	NextId() int
}
