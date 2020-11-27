package store

import "github.com/ultralist/ultralist/ultralist"

// Data is the data structure of what ultralist stores in the .todos.json file.
type Data struct {
	Name     string            `json:"name"`
	Todos    []*ultralist.Todo `json:"todos"`
	Filter   *ultralist.Filter `json:"filter"`
	IsSynced bool              `json:"is_synced"`
}

// Store is the interface for loading and persisting data.
type Store interface {
	Initialize() error
	GetLocation() string
	Load() (*Data, error)
	Save(*Data) error
}
