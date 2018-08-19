package todolist

type Store interface {
	GetLocation() string
	Initialize()
	Load() ([]*Todo, error)
	Save(todos []*Todo)
}
