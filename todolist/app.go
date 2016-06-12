package todolist

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type App struct {
	TodoStore *FileStore
	TodoList  *TodoList
}

func NewApp() *App {
	app := &App{TodoList: &TodoList{}, TodoStore: NewFileStore()}
	return app
}

func (a *App) InitializeRepo() {
	a.TodoStore.Initialize()
}

func (a *App) AddTodo(input string) {
	a.load()
	parser := &Parser{}
	todo := parser.ParseNewTodo(input)

	a.TodoList.Add(todo)
	a.save()
	fmt.Println("Todo added.")
}

func (a *App) DeleteTodo(input string) {
	a.load()
	id := a.getId(input)
	if id != -1 {
		a.TodoList.Delete(id)
		a.save()
		fmt.Println("Todo deleted.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) CompleteTodo(input string) {
	a.load()
	id := a.getId(input)
	if id != -1 {
		a.TodoList.Complete(id)
		a.save()
		fmt.Println("Todo completed.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) UncompleteTodo(input string) {
	a.load()
	id := a.getId(input)
	if id != -1 {
		a.TodoList.Uncomplete(id)
		a.save()
		fmt.Println("Todo uncompleted.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) ArchiveTodo(input string) {
	a.load()
	id := a.getId(input)
	if id != -1 {
		a.TodoList.Archive(id)
		a.save()
		fmt.Println("Todo archived.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) UnarchiveTodo(input string) {
	a.load()
	id := a.getId(input)
	if id != -1 {
		a.TodoList.Unarchive(id)
		a.save()
		fmt.Println("Todo unarchived.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) EditTodoDue(input string) {
	a.load()
	id := a.getId(input)
	if id != -1 {
		todo := a.TodoList.FindById(id)
		parser := &Parser{}
		todo.Due = parser.Due(input, time.Now())
		a.save()
		fmt.Println("Todo due date updated.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) ArchiveCompleted() {
	a.load()
	for _, todo := range a.TodoList.Todos() {
		if todo.Completed {
			todo.Archived = true
		}
	}
	a.save()
	fmt.Println("All archived todos completed.")
}

func (a *App) ListTodos(input string) {
	a.load()
	filtered := NewFilter(a.TodoList.Todos()).Filter(input)
	grouped := a.getGroups(input, filtered)

	formatter := NewFormatter(grouped)
	formatter.Print()
}

func (a *App) getId(input string) int {

	re, _ := regexp.Compile("\\d+")
	if re.MatchString(input) {
		id, _ := strconv.Atoi(re.FindString(input))
		return id
	} else {
		return -1
	}
}

func (a *App) getGroups(input string, todos []*Todo) *GroupedTodos {
	grouper := &Grouper{}
	contextRegex, _ := regexp.Compile(`by c.*$`)
	projectRegex, _ := regexp.Compile(`by p.*$`)

	var grouped *GroupedTodos

	if contextRegex.MatchString(input) {
		grouped = grouper.GroupByContext(todos)
	} else if projectRegex.MatchString(input) {
		grouped = grouper.GroupByProject(todos)
	} else {
		grouped = grouper.GroupByNothing(todos)
	}
	return grouped
}

func (a *App) save() {
	a.TodoStore.Save(a.TodoList.Data)
}
func (a *App) load() {
	a.TodoList.Load(a.TodoStore.Load())
}
