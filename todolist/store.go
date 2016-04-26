package todolist

type Store interface {
	Load()
	Todos() []Todo

	Add(t *Todo)
	Save()

	//Find(id int) Todo
	//Remove(t *Todo)
	NextId() int
}
