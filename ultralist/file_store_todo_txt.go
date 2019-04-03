package ultralist

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

// FileStoreTodoTxt is used to load and store todolists using the todo.txt format.
type FileStoreTodoTxt struct {
	FileLocation string
	Loaded       bool
	TodoList     TodoList
}

// NewFileStoreTodoTxt creates a new FileStoreTodoTxt.
func NewFileStoreTodoTxt() *FileStoreTodoTxt {
	return &FileStoreTodoTxt{FileLocation: "", Loaded: false}
}

// Initialize a brand new file to store todos into.
func (f *FileStoreTodoTxt) Initialize() {
	if f.FileLocation == "" {
		f.FileLocation = f.GetLocation()
	}

	_, err := ioutil.ReadFile(f.FileLocation)
	if err == nil {
		fmt.Println("It looks like a todo.txt file already exists!  Doing nothing.")
		os.Exit(0)
	}

	if err := ioutil.WriteFile(f.FileLocation, []byte("# enter your todos here"), 0644); err != nil {
		fmt.Println("Error writing todo.txt file", err)
	}
}

// Load all todos from a todo.txt file.
func (f *FileStoreTodoTxt) Load() ([]*Todo, error) {
	if f.FileLocation == "" {
		f.FileLocation = f.GetLocation()
	}

	if f.FileLocation == "" {
		fmt.Println("No todo.txt file found!")
		os.Exit(0)
	}

	file, _ := os.Open(f.FileLocation)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var line = scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		f.TodoList.Data = append(f.TodoList.Data, f.ParseLine(line))
	}

	f.Loaded = true

	return f.TodoList.Data, nil
}

// Save all todos to a todo.txt file.
func (f *FileStoreTodoTxt) Save(todoList *TodoList) {
	var data string
	for _, todo := range todoList.Data {
		var line []string

		if todo.Completed {
			line = append(line, "x")
		}

		if todo.IsPriority {
			line = append(line, "(A)")
		}

		if todo.CompletedDate != "" {
			line = append(line, todo.CompletedDate)
		}

		line = append(line, todo.Subject)

		if todo.Due != "" {
			line = append(line, fmt.Sprintf("due:%s", todo.Due))
		}

		line = append(line, fmt.Sprintf("id:%d", todo.ID))

		if todoList.IsSynced {
			line = append(line, fmt.Sprintf("uuid:%s", todo.UUID))
		}

		data += (strings.Join(line, " ") + "\n")
	}
	if err := ioutil.WriteFile(f.FileLocation, []byte(data), 0644); err != nil {
		fmt.Println("Error writing todo.txt file: ", err)
	}
}

// ParseLine parses a todo.txt line into a Todo.
func (f *FileStoreTodoTxt) ParseLine(line string) *Todo {
	todo := NewTodo()
	dateRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

	s := strings.Split(line, " ")

	// find if completed
	if s[0] == "x" {
		todo.Completed = true
		s = append(s[:0], s[1:]...)
	}

	// find if priority
	if s[0] == "(A)" {
		todo.IsPriority = true
		s = append(s[:0], s[1:]...)
	}

	// find a completed date
	if dateRegex.MatchString(s[0]) {
		todo.CompletedDate = s[0]
		s = append(s[:0], s[1:]...)
	}

	// find a created date
	if dateRegex.MatchString(s[0]) {
		// just delete it, we don't store created dates.
		s = append(s[:0], s[1:]...)
	}

	hasDue := false
	hasTodoID := false
	hasUUID := false

	for _, str := range s {
		// find a due date
		if strings.HasPrefix(str, "due:") {
			hasDue = true
			todo.Due = strings.Replace(str, "due:", "", -1)
		}

		if strings.HasPrefix(str, "id:") {
			hasTodoID = true
			idStr := strings.Replace(str, "id:", "", -1)
			id, _ := strconv.Atoi(idStr)
			todo.ID = id
		}

		if strings.HasPrefix(str, "uuid:") {
			hasUUID = true
			todo.UUID = strings.Replace(str, "uuid:", "", -1)
		}
	}

	if hasDue {
		s = s[:len(s)-1]
	}

	if hasTodoID {
		s = s[:len(s)-1]
	} else {
		todo.ID = f.TodoList.NextID()
	}

	if hasUUID {
		s = s[:len(s)-1]
	}

	// all that is left is the subject
	todo.SetSubject(strings.Join(s, " "))

	return todo
}

// GetLocation - Get the location of a todo.txt file.
func (f *FileStoreTodoTxt) GetLocation() string {
	dir, _ := os.Getwd()
	localrepo := fmt.Sprintf("%s/todo.txt", dir)

	usr, _ := user.Current()
	homerepo := fmt.Sprintf("%s/todo.txt", usr.HomeDir)

	_, err := os.Stat(localrepo)
	if err == nil {
		return localrepo
	}

	if _, err = os.Stat(homerepo); err == nil {
		return homerepo
	}

	return ""
}
