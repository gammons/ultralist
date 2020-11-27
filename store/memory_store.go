package store

// MemoryStore defines a store to persist data to Memory, rather than a file.
// This is good for use in tests.
type MemoryStore struct {
	MemoryData *Data
}

// Initialize a new MemoryStore.
func (m *MemoryStore) Initialize() error {
	return nil
}

// Load loads data stored in memory.
func (m *MemoryStore) Load() (*Data, error) {
	return m.MemoryData, nil
}

// GetLocation returns the file location of the MemoryStore.
func (m *MemoryStore) GetLocation() string {
	return ""
}

// Save saves the data to memory.
func (m *MemoryStore) Save(data *Data) error {
	m.MemoryData = data
	return nil
}
