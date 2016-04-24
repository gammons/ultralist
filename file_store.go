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
