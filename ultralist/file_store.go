package ultralist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// FileStore is the main struct of this file.
type FileStore struct {
	FileLocation string
	Loaded       bool
}

// NewFileStore is creating a new file store.
func NewFileStore() *FileStore {
	return &FileStore{FileLocation: "", Loaded: false}
}

// Initialize is initializing a new .todos.json file.
func (f *FileStore) Initialize() {
	if f.FileLocation == "" {
		f.FileLocation = ".todos.json"
	}

	_, err := ioutil.ReadFile(f.FileLocation)
	if err == nil {
		fmt.Println("It looks like a .todos.json file already exists!  Doing nothing.")
		os.Exit(0)
	}
	if err := ioutil.WriteFile(f.FileLocation, []byte("[]"), 0644); err != nil {
		fmt.Println("Error writing json file", err)
		os.Exit(1)
	}
}

// Load is loading a .todos.json file.
func (f *FileStore) Load() ([]*Todo, error) {
	if f.FileLocation == "" {
		f.FileLocation = f.GetLocation()
	}

	data, err := ioutil.ReadFile(f.FileLocation)
	if err != nil {
		fmt.Println("No todo file found!")
		fmt.Println("Initialize a new todo repo by running 'ultralist init'")
		os.Exit(0)
		return nil, err
	}

	var todos []*Todo
	jerr := json.Unmarshal(data, &todos)
	if jerr != nil {
		fmt.Println("Error reading json data", jerr)
		os.Exit(1)
		return nil, jerr
	}
	f.Loaded = true

	return todos, nil
}

// Save is saving a .todos.json file.
func (f *FileStore) Save(todoList *TodoList) {
	// ensure UUID is set for todos at save time
	for _, todo := range todoList.Data {
		if todo.UUID == "" {
			todo.UUID = newUUID()
		}
	}

	data, _ := json.Marshal(todoList.Data)
	if err := ioutil.WriteFile(f.FileLocation, []byte(data), 0644); err != nil {
		fmt.Println("Error writing json file", err)
	}
}

// GetLocation is returning the location of the .todos.json file.
func (f *FileStore) GetLocation() string {
	dir, _ := os.Getwd()
	localrepo := fmt.Sprintf("%s/.todos.json", dir)
	_, ferr := os.Stat(localrepo)
	if ferr == nil {
		return localrepo
	}

	home := UserHomeDir()
	return fmt.Sprintf("%s/.todos.json", home)
}
