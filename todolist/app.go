package todolist

import (
	"fmt"
	"regexp"
	"strings"
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
	a.Load()
	parser := &Parser{}
	todo := parser.ParseNewTodo(input)
	if todo == nil {
		fmt.Println("I need more information. Try something like 'todo a chat with @bob due tom'")
		return
	}

	id := a.TodoList.NextId()
	a.TodoList.Add(todo)
	a.Save()
	fmt.Printf("Todo %d added.\n", id)
}

func (a *App) DeleteTodo(input string) {
	a.Load()
	id, _ := a.getId(input)
	if id == -1 {
		return
	}
	a.TodoList.Delete(id)
	a.Save()
	fmt.Println("Todo deleted.")
}

func (a *App) CompleteTodo(input string) {
	a.Load()
	id, _ := a.getId(input)
	if id == -1 {
		return
	}
	a.TodoList.Complete(id)
	a.Save()
	fmt.Println("Todo completed.")
}

func (a *App) UncompleteTodo(input string) {
	a.Load()
	id, _ := a.getId(input)
	if id == -1 {
		return
	}
	a.TodoList.Uncomplete(id)
	a.Save()
	fmt.Println("Todo uncompleted.")
}

func (a *App) ArchiveTodo(input string) {
	a.Load()
	id, _ := a.getId(input)
	if id == -1 {
		return
	}
	a.TodoList.Archive(id)
	a.Save()
	fmt.Println("Todo archived.")
}

func (a *App) UnarchiveTodo(input string) {
	a.Load()
	id, _ := a.getId(input)
	if id == -1 {
		return
	}
	a.TodoList.Unarchive(id)
	a.Save()
	fmt.Println("Todo unarchived.")
}

func (a *App) EditTodoSubject(input string) {
	a.Load()
	_, id, subject := Parser{input}.Parse()
	id, todo := a.getId(input)
	if id == -1 {
		return
	}
	todo.Subject = subject
	a.Save()
	fmt.Println("Todo subject updated.")
}

func (a *App) EditTodoDue(input string) {
	a.Load()
	id, todo := a.getId(input)
	if id == -1 {
		return
	}
	parser := &Parser{}
	todo.Due = parser.Due(input, time.Now())
	a.Save()
	fmt.Println("Todo due date updated.")
}

func (a *App) ExpandTodo(input string) {
	a.Load()
	id, _ := a.getId(input)
	parser := &Parser{}
	if id == -1 {
		return
	}

	commonProject := parser.ExpandProject(input)
	todos := strings.LastIndex(input, ":")
	if commonProject == "" || len(input) <= todos+1 || todos == -1 {
		fmt.Println("I'm expecting a format like \"todolist ex <project>: <todo1>, <todo2>, ...\"")
		return
	}

	newTodos := strings.Split(input[todos+1:], ",")

	for _, todo := range newTodos {
		args := []string{"add ", commonProject, " ", todo}
		a.AddTodo(strings.Join(args, ""))
	}

	a.TodoList.Delete(id)
	a.Save()
	fmt.Println("Todo expanded.")
}

func (a *App) ArchiveCompleted() {
	a.Load()
	for _, todo := range a.TodoList.Todos() {
		if todo.Completed {
			todo.Archived = true
		}
	}
	a.Save()
	fmt.Println("All completed todos have been archived.")
}

func (a *App) ListTodos(input string) {
	a.Load()
	filtered := NewFilter(a.TodoList.Todos()).Filter(input)
	grouped := a.getGroups(input, filtered)

	formatter := NewFormatter(grouped)
	formatter.Print()
}

func (a *App) PrioritizeTodo(input string) {
	a.Load()
	id, _ := a.getId(input)
	if id == -1 {
		return
	}
	a.TodoList.Prioritize(id)
	a.Save()
	fmt.Println("Todo prioritized.")
}

func (a *App) UnprioritizeTodo(input string) {
	a.Load()
	id, _ := a.getId(input)
	if id == -1 {
		return
	}
	a.TodoList.Unprioritize(id)
	a.Save()
	fmt.Println("Todo un-prioritized.")
}

func (a *App) getId(input string) (int, *Todo) {
	_, id, _ := Parser{input}.Parse()
	todo := a.TodoList.FindById(id)
	if todo == nil {
		fmt.Println("No such id.")
		return -1, nil
	}
	return id, todo
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

func (a *App) GarbageCollect() {
	a.Load()
	a.TodoList.GarbageCollect()
	a.Save()
	fmt.Println("Garbage collection complete.")
}

func (a *App) Load() error {
	todos, err := a.TodoStore.Load()
	if err != nil {
		return err
	}
	a.TodoList.Load(todos)
	return nil
}

func (a *App) Save() {
	a.TodoStore.Save(a.TodoList.Data)
}
