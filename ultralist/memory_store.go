package ultralist

// MemoryStore is the main struct of this file.
type MemoryStore struct {
	Todos []*Todo
}

// NewMemoryStore is starting new memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

// Initialize is initializing a new memory store.
func (m *MemoryStore) Initialize() {}

// Load is loading todos from the memory store.
func (m *MemoryStore) Load() ([]*Todo, error) {
	return m.Todos, nil
}

func (m *MemoryStore) LocalTodosFileExists() bool {
	return false
}

// Save is saving todos to the memory store.
func (m *MemoryStore) Save(todos []*Todo) {
	m.Todos = todos
}

// GetLocation is giving the location of the memory store.
func (m *MemoryStore) GetLocation() string {
	return ""
}
