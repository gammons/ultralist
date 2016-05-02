package todolist

import (
	"fmt"
	"regexp"
	"strconv"
)

type App struct {
	TodoStore Store
}

func NewApp() *App {
	app := &App{TodoStore: NewFileStore()}
	app.TodoStore.Load()
	return app
}

func (a *App) AddTodo(input string) {
	parser := &Parser{}
	todo := parser.ParseNewTodo(input)

	a.TodoStore.Add(todo)
	a.TodoStore.Save()
	fmt.Println("Todo added.")
}

func (a *App) DeleteTodo(input string) {
	id := a.getId(input)
	if id != -1 {
		a.TodoStore.Delete(id)
		a.TodoStore.Save()
		fmt.Println("Todo deleted.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) CompleteTodo(input string) {
	id := a.getId(input)
	if id != -1 {
		a.TodoStore.Complete(id)
		a.TodoStore.Save()
		fmt.Println("Todo completed.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) UncompleteTodo(input string) {
	id := a.getId(input)
	if id != -1 {
		a.TodoStore.Uncomplete(id)
		a.TodoStore.Save()
		fmt.Println("Todo uncompleted.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) ArchiveTodo(input string) {
	id := a.getId(input)
	if id != -1 {
		a.TodoStore.Archive(id)
		a.TodoStore.Save()
		fmt.Println("Todo archived.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) UnarchiveTodo(input string) {
	id := a.getId(input)
	if id != -1 {
		a.TodoStore.Unarchive(id)
		a.TodoStore.Save()
		fmt.Println("Todo unarchived.")
	} else {
		fmt.Println("Could not find id.")
	}
}

func (a *App) ListTodos(input string) {
	filtered := NewFilter(a.TodoStore.Todos()).Filter(input)
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

func (a *App) getGroups(input string, todos []Todo) *GroupedTodos {
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
