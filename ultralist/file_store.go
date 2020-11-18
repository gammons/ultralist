package ultralist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TodosJSONFile is the filename to store todos in
const TodosJSONFile = ".todos.json"

// FileStore is the main struct of this file.
type FileStore struct {
	Loaded bool
}

// NewFileStore is creating a new file store.
func NewFileStore() *FileStore {
	return &FileStore{Loaded: false}
}

// Initialize is initializing a new .todos.json file.
func (f *FileStore) Initialize() {
	if f.LocalTodosFileExists() {
		fmt.Println("It looks like a .todos.json file already exists!  Doing nothing.")
		os.Exit(0)
	}
	if err := ioutil.WriteFile(TodosJSONFile, []byte("[]"), 0644); err != nil {
		fmt.Println("Error writing json file", err)
		os.Exit(1)
	}
}

// Returns if a local .todos.json file exists in the current dir.
func (f *FileStore) LocalTodosFileExists() bool {
	dir, _ := os.Getwd()
	localrepo := filepath.Join(dir, TodosJSONFile)
	_, err := os.Stat(localrepo)
	return err == nil
}

// Load is loading a .todos.json file, either from cwd, or the home directory
func (f *FileStore) Load() ([]*Todo, error) {
	data, err := ioutil.ReadFile(f.GetLocation())
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
func (f *FileStore) Save(todos []*Todo) {
	// ensure UUID is set for todos at save time
	for _, todo := range todos {
		if todo.UUID == "" {
			todo.UUID = newUUID()
		}
	}

	data, _ := json.Marshal(todos)
	if err := ioutil.WriteFile(f.GetLocation(), []byte(data), 0644); err != nil {
		fmt.Println("Error writing json file", err)
	}
}

// GetLocation is returning the location of the .todos.json file.
func (f *FileStore) GetLocation() string {
	if f.LocalTodosFileExists() {
		dir, _ := os.Getwd()
		localrepo := filepath.Join(dir, TodosJSONFile)
		return localrepo
	}
	return fmt.Sprintf("%s/%s", UserHomeDir(), TodosJSONFile)
}
