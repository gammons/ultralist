package todolist

type Store interface {
	Load()
	Todos() []Todo

	Add(t *Todo)
	Save()

	IndexOf(t *Todo) int
	FindById(id int) *Todo
	Delete(id int)
	NextId() int
}
