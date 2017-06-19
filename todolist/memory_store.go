package todolist

type MemoryStore struct {
	Todos []*Todo
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Initialize() {}

func (m *MemoryStore) Load() ([]*Todo, error) {
	return m.Todos, nil
}

func (m *MemoryStore) Save(todos []*Todo) {
	m.Todos = todos
}
