package todolist

type Store interface {
	Initialize()
	Load() ([]*Todo, error)
	Save(todos []*Todo)
}
