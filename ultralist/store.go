package ultralist

// Store is the interface for ultralist todos.
type Store interface {
	GetLocation() string
	LocalTodosFileExists() bool
	Initialize()
	Load() ([]*Todo, error)
	Save(todos []*Todo)
}
