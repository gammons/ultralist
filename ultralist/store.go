package ultralist

// Store is the interface for ultralist todos.
type Store interface {
	GetLocation() string
	Initialize()
	Load() ([]*Todo, error)
	Save(todos []*Todo)
}
