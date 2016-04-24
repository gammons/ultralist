package todolist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type FileStore struct {
	FileLocation string
	Data         []Todo
}

func NewFileStore() *FileStore {
	return &FileStore{FileLocation: "todos.json"}
}

func (f *FileStore) Load() {
	data, err := ioutil.ReadFile(f.FileLocation)
	if err != nil {
		fmt.Println("Error reading file", err)
	}

	jerr := json.Unmarshal(data, &f.Data)
	if jerr != nil {
		fmt.Println("Error reading json data", jerr)
	}

}
