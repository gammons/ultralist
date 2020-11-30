package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ultralist/ultralist/ultralist"
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
func (f *FileStore) Initialize() error {
	if f.localTodosFileExists() {
		return errors.New("It looks like a .todos.json file already exists")
	}

	marshalled, _ := json.Marshal(&Data{})
	if err := ioutil.WriteFile(TodosJSONFile, []byte(marshalled), 0644); err != nil {
		return err
	}

	return nil
}

// Load is loading a .todos.json file, either from cwd, or the home directory
func (f *FileStore) Load() (*Data, error) {
	fileData, err := ioutil.ReadFile(f.GetLocation())
	if err != nil {
		return nil, err
	}

	var data *Data
	err = json.Unmarshal(fileData, &data)

	if err != nil {

		// try loading the .todos.json legacy file format before giving up.
		data, err = f.parseLegacyTodosJSONFile(fileData)

		if err != nil {
			return nil, err
		}
	}
	f.Loaded = true

	return data, nil
}

// Save will persist the Data into a .todos.json file.
func (f *FileStore) Save(data *Data) error {
	marshalled, _ := json.Marshal(data)
	if err := ioutil.WriteFile(f.GetLocation(), []byte(marshalled), 0644); err != nil {
		return err
	}
	return nil
}

func (f *FileStore) parseLegacyTodosJSONFile(fileData []byte) (*Data, error) {
	var todos []*ultralist.Todo
	err := json.Unmarshal(fileData, &todos)
	if err != nil {
		return nil, err
	}

	return &Data{
		TodoList: &ultralist.TodoList{Data: todos},
	}, nil
}

// Returns if a local .todos.json file exists in the current dir.
func (f *FileStore) localTodosFileExists() bool {
	dir, _ := os.Getwd()
	localrepo := filepath.Join(dir, TodosJSONFile)
	_, err := os.Stat(localrepo)
	return err == nil
}

// GetLocation is returning the location of the .todos.json file.
func (f *FileStore) GetLocation() string {
	if f.localTodosFileExists() {
		dir, _ := os.Getwd()
		localrepo := filepath.Join(dir, TodosJSONFile)
		return localrepo
	}
	return fmt.Sprintf("%s/%s", f.userHomeDir(), TodosJSONFile)
}

// UserHomeDir returns the home dir of the current user.
func (f *FileStore) userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}
