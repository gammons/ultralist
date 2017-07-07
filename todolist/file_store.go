package todolist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lmcm.io/stddirs"
	"os"
	"os/user"
	"path/filepath"
)

const APP_ID = "com.github.gammons.todolist"

type FileStore struct {
	FileLocation string
	Loaded       bool
}

func NewFileStore() *FileStore {
	return &FileStore{FileLocation: "", Loaded: false}
}

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
	fmt.Println("Todo repo initialized.")
}

func (f *FileStore) Load() ([]*Todo, error) {
	if f.FileLocation == "" {
		var err error
		f.FileLocation, err = getLocation()
		if err != nil {
			fmt.Println("Error whlie finding todo file", err)
		}
	}

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

	return todos, nil
}

func (f *FileStore) Save(todos []*Todo) {
	os.MkdirAll(filepath.Dir(f.FileLocation), 0700)

	data, _ := json.Marshal(todos)
	if err := ioutil.WriteFile(f.FileLocation, []byte(data), 0644); err != nil {
		fmt.Println("Error writing json file", err)
	}
}

func getLocation() (location string, err error) {
	// Look for `.todos.json` in current working directory
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	localrepo := filepath.Join(wd, ".todos.json")
	if fileExists(localrepo) {
		return localrepo, nil
	}

	// Look for `.todos.json` in user's home directory
	// (For compatability with old inits, see #72)
	usr, err := user.Current()
	if err != nil {
		return
	}
	homerepo := filepath.Join(usr.HomeDir, ".todos.json")
	if fileExists(homerepo) {
		return homerepo, nil
	}

	// Default to the global todolist, whether it exists or not
	globalrepo := filepath.Join(stddirs.DataDir(APP_ID), "todos.json")
	return globalrepo, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
