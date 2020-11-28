package store

import "github.com/ultralist/ultralist/ultralist"

// Data is the data structure of what ultralist stores in the .todos.json file.
type Data struct {
	TodoList *ultralist.TodoList `json:"todo_list"`
	Filter   *ultralist.Filter   `json:"filter"`
}

// Store is the interface for loading and persisting data.
type Store interface {
	Initialize() error
	GetLocation() string
	Load() (*Data, error)
	Save(*Data) error
}
