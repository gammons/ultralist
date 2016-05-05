package todolist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"sort"
)

type FileStore struct {
	FileLocation string
	Data         []*Todo
}

func NewFileStore() *FileStore {
	usr, _ := user.Current()
	return &FileStore{FileLocation: usr.HomeDir + "/.todos.json"}
}

func (f *FileStore) Add(todo *Todo) {
	todo.Id = f.NextId()
	f.Data = append(f.Data, todo)
}

func (f *FileStore) FindById(id int) *Todo {
	for _, todo := range f.Data {
		if todo.Id == id {
			return todo
		}
	}
	return nil
}

func (f *FileStore) Delete(id int) {
	i := -1
	for index, todo := range f.Data {
		if todo.Id == id {
			i = index
		}
	}

	f.Data = append(f.Data[:i], f.Data[i+1:]...)
}

func (f *FileStore) Complete(id int) {
	todo := f.FindById(id)
	todo.Completed = true
	f.Delete(id)
	f.Data = append(f.Data, todo)
}

func (f *FileStore) Uncomplete(id int) {
	todo := f.FindById(id)
	todo.Completed = false
	f.Delete(id)
	f.Data = append(f.Data, todo)
}

func (f *FileStore) Archive(id int) {
	todo := f.FindById(id)
	todo.Archived = true
	f.Delete(id)
	f.Data = append(f.Data, todo)
}

func (f *FileStore) Unarchive(id int) {
	todo := f.FindById(id)
	todo.Archived = false
	f.Delete(id)
	f.Data = append(f.Data, todo)
}

func (f *FileStore) IndexOf(todoToFind *Todo) int {
	for i, todo := range f.Data {
		if todo.Id == todoToFind.Id {
			return i
		}
	}
	return -1
}

func (f *FileStore) Load() {
	data, err := ioutil.ReadFile(f.FileLocation)
	if err != nil {
		fmt.Println("No todo file found at ~/.todos.json!")
		fmt.Println("Initializing new todo repo at ~/.todos.json")
		f.Initialize()
		return
	}

	jerr := json.Unmarshal(data, &f.Data)
	if jerr != nil {
		fmt.Println("Error reading json data", jerr)
		os.Exit(1)
	}
}

func (f *FileStore) Initialize() {
	if err := ioutil.WriteFile(f.FileLocation, []byte("[]"), 0644); err != nil {
		fmt.Println("Error writing json file", err)
	}
}

func (f *FileStore) Save() {
	data, _ := json.Marshal(f.Data)
	if err := ioutil.WriteFile(f.FileLocation, []byte(data), 0644); err != nil {
		fmt.Println("Error writing json file", err)
	}
}

type ByDate []*Todo

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	t1Due := a[i].CalculateDueTime()
	t2Due := a[j].CalculateDueTime()
	return t1Due.Before(t2Due)
}

func (f *FileStore) Todos() []*Todo {
	sort.Sort(ByDate(f.Data))
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
