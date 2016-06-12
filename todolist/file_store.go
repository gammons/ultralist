package todolist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type FileStore struct {
	FileLocation string
	Loaded       bool
}

func NewFileStore() *FileStore {
	return &FileStore{FileLocation: ".todos.json", Loaded: false}
}

func (f *FileStore) Load() []*Todo {
	data, err := ioutil.ReadFile(f.FileLocation)
	if err != nil {
		fmt.Println("No todo file found!")
		fmt.Println("Initialize a new todo repo by running 'todo init'")
		os.Exit(0)
	}

	var todos []*Todo
	jerr := json.Unmarshal(data, &todos)
	if jerr != nil {
		fmt.Println("Error reading json data", jerr)
		os.Exit(1)
	}
	f.Loaded = true

	return todos
}

func (f *FileStore) Initialize() {
	_, err := ioutil.ReadFile(f.FileLocation)
	if err == nil {
		fmt.Println("It looks like a .todos.json file already exists!  Doing nothing.")
		os.Exit(0)
	}
	if err := ioutil.WriteFile(f.FileLocation, []byte("[]"), 0644); err != nil {
		fmt.Println("Error writing json file", err)
		os.Exit(1)
	}
	fmt.Println("Todo repo initialized.")
}

func (f *FileStore) Save(todos []*Todo) {
	data, _ := json.Marshal(todos)
	if err := ioutil.WriteFile(f.FileLocation, []byte(data), 0644); err != nil {
		fmt.Println("Error writing json file", err)
	}
}
