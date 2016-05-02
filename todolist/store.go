package todolist

type Store interface {
	Load()
	Save()
	Todos() []Todo

	Add(t *Todo)
	Delete(id int)

	Complete(id int)
	Uncomplete(id int)

	Archive(id int)
	Unarchive(id int)

	IndexOf(t *Todo) int
	FindById(id int) *Todo
	NextId() int
}
