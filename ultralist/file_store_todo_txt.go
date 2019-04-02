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

type FileStoreTodoTxt struct {
	FileLocation string
	Loaded       bool
}

func NewFileStoreTodoTxt() *FileStoreTodoTxt {
	return &FileStoreTodoTxt{FileLocation: "", Loaded: false}
}

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

func (f *FileStoreTodoTxt) HasTodoTxtFile() bool {
	return f.GetLocation() != ""
}

func (f *FileStoreTodoTxt) Load() ([]*Todo, error) {
	if f.FileLocation == "" {
		f.FileLocation = f.GetLocation()
	}

	if f.FileLocation == "" {
		fmt.Println("No todo.txt file found!")
		os.Exit(0)
	}

	var todos []*Todo

	file, _ := os.Open(f.FileLocation)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var line = scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		fmt.Println("line: ", line)
		todos = append(todos, f.ParseLine(line))
	}

	f.Loaded = true

	return todos, nil
}

func (f *FileStoreTodoTxt) ParseLine(line string) *Todo {
	todo := &Todo{}
	dateRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	dueRegex := regexp.MustCompile(`due:\d{4}-\d{2}-\d{2}`)
	todoIDRegex := regexp.MustCompile(`id:\d+`)

	s := strings.Split(line, " ")

	// find if completed
	if s[0] == "x" {
		fmt.Println("this would be completed")
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

	for i, str := range s {
		// find a due date
		if dueRegex.MatchString(str) {
			hasDue = true
			todo.Due = strings.Replace(s[i], "due:", "", -1)
		}
		if todoIDRegex.MatchString(str) {
			hasTodoID = true
			idStr := strings.Replace(s[i], "id:", "", -1)
			id, _ := strconv.Atoi(idStr)
			fmt.Println("idStr = ", idStr)
			todo.ID = id
		}
	}

	if hasDue {
		s = s[:len(s)-1]
	}

	if hasTodoID {
		s = s[:len(s)-1]
	}

	// all that is left is the subject
	todo.Subject = strings.Join(s, " ")

	return todo
}

func (f *FileStoreTodoTxt) GetLocation() string {
	dir, _ := os.Getwd()
	fmt.Println("dir is ", dir)
	localrepo := fmt.Sprintf("%s/todo.txt", dir)

	usr, _ := user.Current()
	homerepo := fmt.Sprintf("%s/todo.txt", usr.HomeDir)

	_, err := os.Stat(localrepo)
	if err == nil {
		return localrepo
	}

	if _, err = os.Stat(localrepo); err != nil {
		return homerepo
	}

	return ""
}
