package todolist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

type FileStore struct {
	FileLocation string
	Data         []Todo
}

func NewFileStore() *FileStore {
	usr, _ := user.Current()
	return &FileStore{FileLocation: usr.HomeDir + "/.todos.json"}
}

func (f *FileStore) Add(todo *Todo) {
	todo.Id = f.NextId()
	f.Data = append(f.Data, *todo)
}

func (f *FileStore) Load() {
	data, err := ioutil.ReadFile(f.FileLocation)
	if err != nil {
		fmt.Println("Error reading file", err)
		os.Exit(1)
	}

	jerr := json.Unmarshal(data, &f.Data)
	if jerr != nil {
		fmt.Println("Error reading json data", jerr)
		os.Exit(1)
	}
}

func (f *FileStore) Save() {
	data, _ := json.Marshal(f.Data)
	if err := ioutil.WriteFile(f.FileLocation, []byte(data), 0644); err != nil {
		fmt.Println("Error writing json file", err)
	}
}

func (f *FileStore) Todos() []Todo {
	return f.Data
}

func (f *FileStore) NextId() int {
	maxId := 0
	for _, todo := range f.Data {
		if todo.Id > maxId {
			maxId = todo.Id
		}
	}
	return maxId + 1
}
