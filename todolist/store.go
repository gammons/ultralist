package todolist

type Store interface {
	Load()
	Save()

	Find(id int) Todo
	Add(t *Todo)
	Remove(t *Todo)
	NextId() int
}
