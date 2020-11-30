package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/skratchdot/open-golang/open"
	"github.com/ultralist/ultralist/store"
	"github.com/ultralist/ultralist/sync"
	"github.com/ultralist/ultralist/tui"
	"github.com/ultralist/ultralist/ultralist"
)

// App is a representation of the ultralist app that is invoked in CLI mode.
// it will output to stdout.
type App struct {
	TodoStore store.Store
	TodoList  *ultralist.TodoList
	Filter    *ultralist.Filter
}

func NewApp() *App {
	return &App{
		TodoStore: store.NewFileStore(),
		Filter:    &ultralist.Filter{},
	}
}

// InitializeRepo will initialize a new .todos.json repo and then tell the user.
func (a *App) InitializeRepo() {
	a.TodoStore.Initialize()

	fmt.Println("Repo initialized.")
}

// AddTodo adds a todo to the current todolist via the CLI
// Takes a string `input`, runs it through Ultralist's InputParser, and then adds the todo to the list
func (a *App) AddTodo(input string) {
	parser := &InputParser{}

	filter, _, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("I need more information. Try something like 'ultralist a chat with @bob due tom'")
		return
	}

	todoItem, err := ultralist.CreateTodo(filter)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	a.loadTodoList()
	a.TodoList.Add(todoItem)
	a.saveTodoList()

	fmt.Printf("Todo %d added.\n", todoItem.ID)
}

// ArchiveCompletedTodos will archive all completed todos.
func (a *App) ArchiveCompletedTodos() {
	a.loadTodoList()
	a.TodoList.ArchiveCompletedTodos()
	a.saveTodoList()
}

// ArchiveTodo will archive todos with the specified ids.
func (a *App) ArchiveTodos(ids ...int) {
	a.loadTodoList()
	a.TodoList.Archive(ids...)
	a.saveTodoList()

	fmt.Printf("%s archived.\n", a.pluralize("Todo", len(ids)))
}

// AuthWorkflow starts the procedure for authenticating against ultralist.io.
func (a *App) AuthWorkflow() {
	webapp := &sync.Webapp{}
	backend := sync.NewBackend()

	open.Start(backend.AuthURL())
	fmt.Println("Head to your browser to complete authorization steps.")
	fmt.Println("Listening for response...")
	webapp.Run()
}

// CheckAuth will check a user's authentication against the ultralist.io service.
func (a *App) CheckAuth() {
	synchronizer := sync.NewSynchronizer()
	name, err := synchronizer.CheckAuth()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Hello, %s! Your credentials are valid.", name)
}

// CompleteTodos will complete todos with the specified ids.
func (a *App) CompleteTodos(ids ...int) {
	a.loadTodoList()
	a.TodoList.Complete(ids...)
	a.saveTodoList()

	fmt.Printf("%s completed.\n", a.pluralize("Todo", len(ids)))
}

// DeleteTodos will complete todos with the specified ids.
func (a *App) DeleteTodos(ids ...int) {
	a.loadTodoList()
	a.TodoList.Delete(ids...)
	a.saveTodoList()

	fmt.Printf("%s deleted.\n", a.pluralize("Todo", len(ids)))
}

// EditTodo will edit a todo.
func (a *App) EditTodo(todoID int, input string) {
	a.loadTodoList()

	parser := &InputParser{}

	filter, _, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("I need more information. Try something like 'ultralist a chat with @bob due tom'")
		return
	}

	todo := a.TodoList.FindByID(todoID)
	if todo == nil {
		fmt.Println("No todo with that id.")
		return
	}

	if err = ultralist.EditTodo(todo, a.TodoList, filter); err != nil {
		fmt.Println(err.Error())
		return
	}

	a.saveTodoList()

	fmt.Println("Todo updated.")
}

// ListTodos will list all todos with the specified options.
func (a *App) ListTodos(input string, printer Printer) {
	a.loadTodoList()

	parser := &InputParser{}

	filter, grouping, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	todoFilter := &ultralist.TodoFilter{Todos: a.TodoList.Todos(), Filter: filter}
	grouper := &ultralist.Grouper{}
	groups := grouper.GroupTodos(todoFilter.ApplyFilter(), grouping)
	printer.Print(groups)
}

// UncompleteTodos will complete todos with the specified ids.
func (a *App) UncompleteTodos(ids ...int) {
	a.loadTodoList()
	a.TodoList.Complete(ids...)
	a.saveTodoList()

	fmt.Printf("%s uncompleted.\n", a.pluralize("Todo", len(ids)))
}

// PrioritizeTodos will prioritize todos with the specified IDs.
func (a *App) PrioritizeTodos(ids ...int) {
	a.loadTodoList()
	a.TodoList.Prioritize(ids...)
	a.saveTodoList()

	fmt.Printf("%s prioritized.\n", a.pluralize("Todo", len(ids)))
}

// GarbageCollect will delete all archived todos, thus re-claiming
// low todo IDs.
func (a *App) GarbageCollect() {
	a.loadTodoList()
	a.TodoList.GarbageCollect()
	a.saveTodoList()
	fmt.Println("Garbage collection complete.")
}

// OpenManager opens the tui manager
func (a *App) OpenManager() {
	a.loadTodoList()
	manager := tui.NewManager(a.TodoList)
	manager.RunManager()
}

// SetTodosStatus will set the status for the specified todo ids
func (a *App) SetTodosStatus(status string, ids ...int) {
	a.loadTodoList()
	a.TodoList.SetStatus(status, ids...)
	a.saveTodoList()

	fmt.Printf("%s status set.\n", a.pluralize("Todo", len(ids)))
}

// SetupSync will run through the process of syncing a list with ultralist.io.
func (a *App) SetupSync() {
	backend := sync.NewBackend()
	if !backend.CredsFileExists() {
		fmt.Println("You're not authenticated with ultralist.io yet.  Please run `ultralist auth` first.")
		return
	}

	if _, err := os.Stat(a.TodoStore.GetLocation()); err == nil {
		a.setupSyncForExistingList(backend)
		return
	}

	// at this point, it is known that a local todos file does not exist.
	type Response struct {
		Todolists []ultralist.TodoList `json:"todolists"`
	}

	var response *Response

	resp, err := backend.PerformRequest("GET", "/api/v1/todo_lists", []byte{})
	if err != nil {
		fmt.Println(err)
		return
	}

	json.Unmarshal(resp, &response)

	var todolistNames []string
	for _, todolist := range response.Todolists {
		todolistNames = append(todolistNames, todolist.Name)
	}
	prompt2 := promptui.Select{
		Label: "Choose a list to import",
		Items: todolistNames,
	}

	idx, _, _ := prompt2.Run()
	a.TodoList = &response.Todolists[idx]
	a.TodoStore.Initialize()
	a.TodoStore.Save(&store.Data{TodoList: a.TodoList})

	// a.EventLogger = NewEventLogger(a.TodoList, a.TodoStore)
	// a.EventLogger.WriteSyncedLists()
}

// UnarchiveTodos will unarchive todos with the specified IDs.
func (a *App) UnarchiveTodos(ids ...int) {
	a.loadTodoList()
	a.TodoList.Unarchive(ids...)
	a.saveTodoList()

	fmt.Printf("%s unarchived.\n", a.pluralize("Todo", len(ids)))
}

// UnprioritizeTodos will un-prioritize todos with the specified IDs.
func (a *App) UnprioritizeTodos(ids ...int) {
	a.loadTodoList()
	a.TodoList.Unprioritize(ids...)
	a.saveTodoList()

	fmt.Printf("%s unprioritized.\n", a.pluralize("Todo", len(ids)))
}

func (a *App) loadTodoList() {
	var data *store.Data
	data, err := a.TodoStore.Load()

	if err != nil {
		fmt.Println("I had trouble loading the .todos.json file.")
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}

	a.TodoList = data.TodoList
}

func (a *App) saveTodoList() {
	data := &store.Data{
		TodoList: a.TodoList,
		Filter:   &ultralist.Filter{},
	}

	if err := a.TodoStore.Save(data); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (a *App) setupSyncForExistingList(backend *sync.Backend) {
	a.loadTodoList()

	if a.TodoList.IsSynced {
		fmt.Println("This list is already sycned with ultralist.io. Use the --unsync flag to stop syncing this list.")
		return
	}

	prompt := promptui.Select{
		Label: "You have a todos list in this directory.  What would you like to do?",
		Items: []string{"Sync my list to ultralist.io", "Pull a list from ultralist.io, replacing the list that's here"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return
	}
	if strings.HasPrefix(result, "Sync my list") {
		prompt := promptui.Prompt{
			Label: "Give this list a name",
		}

		result, err := prompt.Run()
		if err != nil {
			fmt.Println("A name is required to sync a list.")
			return
		}
		// a.EventLogger.CurrentSyncedList.Name = result
		a.TodoList.Name = result
		a.TodoList.IsSynced = true
		// a.EventLogger.WriteSyncedLists()

		// create a todo list via the API.
		backend.CreateTodoList(a.TodoList)
		return
	}
}

func (a *App) pluralize(name string, length int) string {
	if length == 1 {
		return name
	}
	return name + "s"
}
