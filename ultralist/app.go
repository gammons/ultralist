package ultralist

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/skratchdot/open-golang/open"
)

// Current version of ultralist.
const (
	VERSION string = "0.9.3"
)

// App is the giving you the structure of the ultralist app.
type App struct {
	EventLogger *EventLogger
	TodoStore   Store
	Printer     Printer
	TodoList    *TodoList
}

// NewApp is creating a new ultralist app.
func NewApp() *App {

	app := &App{
		TodoList: &TodoList{},
		Printer:  NewScreenPrinter(),
	}

	storeTodoTxt := NewFileStoreTodoTxt()

	if storeTodoTxt.GetLocation() != "" {
		app.TodoStore = storeTodoTxt
	} else {
		app.TodoStore = NewFileStore()
	}

	return app
}

// InitializeRepo is initializing ultralist repo.
func (a *App) InitializeRepo() {
	a.TodoStore.Initialize()
	fmt.Println("Repo initialized.")

	backend := NewBackend()
	eventLogger := &EventLogger{Store: a.TodoStore}
	eventLogger.LoadSyncedLists()

	if !backend.CredsFileExists() {
		return
	}

	prompt := promptui.Prompt{
		Label:     "Do you wish to sync this list with ultralist.io",
		IsConfirm: true,
	}

	result, _ := prompt.Run()
	if result != "y" {
		return
	}

	if !backend.CanConnect() {
		fmt.Println("I cannot connect to ultralist.io right now.")
		return
	}

	// fetch lists from ultralist.io, or allow user to create a new list
	// use the "select_add" example in promptui as a way to do this
	type Response struct {
		Todolists []TodoList `json:"todolists"`
	}

	var response *Response

	resp := backend.PerformRequest("GET", "/api/v1/todo_lists", []byte{})
	json.Unmarshal(resp, &response)

	var todolistNames []string
	for _, todolist := range response.Todolists {
		todolistNames = append(todolistNames, todolist.Name)
	}

	prompt2 := promptui.SelectWithAdd{
		Label:    "You can sync with an existing list on ultralist, or create a new list.",
		Items:    todolistNames,
		AddLabel: "New list...",
	}

	idx, name, _ := prompt2.Run()
	if idx == -1 {
		eventLogger.CurrentSyncedList.Name = name
	} else {
		eventLogger.CurrentSyncedList.Name = response.Todolists[idx].Name
		eventLogger.CurrentSyncedList.UUID = response.Todolists[idx].UUID
		a.TodoList = &response.Todolists[idx]
		a.save()
	}

	eventLogger.WriteSyncedLists()
}

// AddTodo is adding a new todo.
func (a *App) AddTodo(input string) {
	a.Load()
	parser := &Parser{}
	todo := parser.ParseNewTodo(input)
	if todo == nil {
		fmt.Println("I need more information. Try something like 'todo a chat with @bob due tom'")
		return
	}

	id := a.TodoList.NextID()
	a.TodoList.Add(todo)
	a.save()
	fmt.Printf("Todo %d added.\n", id)
}

// AddDoneTodo adds a new todo and immediately completed it.
func (a *App) AddDoneTodo(input string) {
	a.Load()

	r, _ := regexp.Compile(`^(done)(\s*|)`)
	input = r.ReplaceAllString(input, "")
	parser := &Parser{}
	todo := parser.ParseNewTodo(input)
	if todo == nil {
		fmt.Println("I need more information. Try something like 'todo done chating with @bob'")
		return
	}

	id := a.TodoList.NextID()
	a.TodoList.Add(todo)
	a.TodoList.Complete(id)
	a.save()
	fmt.Printf("Completed Todo %d added.\n", id)
}

// DeleteTodo deletes a todo.
func (a *App) DeleteTodo(input string) {
	a.Load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Delete(ids...)
	a.save()
	fmt.Printf("%s deleted.\n", pluralize(len(ids), "Todo", "Todos"))
}

// CompleteTodo completes a todo.
func (a *App) CompleteTodo(input string, archive bool) {
	a.Load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Complete(ids...)
	if archive {
		a.TodoList.Archive(ids...)
	}
	a.save()
	fmt.Println("Todo completed.")
}

// UncompleteTodo uncompletes a todo.
func (a *App) UncompleteTodo(input string) {
	a.Load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Uncomplete(ids...)
	a.save()
	fmt.Println("Todo uncompleted.")
}

// ArchiveTodo archives a todo.
func (a *App) ArchiveTodo(input string) {
	a.Load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Archive(ids...)
	a.save()
	fmt.Println("Todo archived.")
}

// UnarchiveTodo unarchives a todo.
func (a *App) UnarchiveTodo(input string) {
	a.Load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Unarchive(ids...)
	a.save()
	fmt.Println("Todo unarchived.")
}

// EditTodo edits a todo with the given input.
func (a *App) EditTodo(input string) {
	a.Load()
	id := a.getID(input)
	if id == -1 {
		return
	}
	todo := a.TodoList.FindByID(id)
	if todo == nil {
		fmt.Println("No such id.")
		return
	}
	parser := &Parser{}

	if parser.ParseEditTodo(todo, input) {
		a.save()
		fmt.Println("Todo updated.")
	}
}

// ExpandTodo expands a todo.
func (a *App) ExpandTodo(input string) {
	a.Load()
	id := a.getID(input)
	parser := &Parser{}
	if id == -1 {
		return
	}

	commonProject := parser.ExpandProject(input)
	todos := strings.LastIndex(input, ":")
	if commonProject == "" || len(input) <= todos+1 || todos == -1 {
		fmt.Println("I'm expecting a format like \"ultralist expand <project>: <todo1>, <todo2>, ...")
		return
	}

	newTodos := strings.Split(input[todos+1:], ",")

	for _, todo := range newTodos {
		args := []string{"add ", commonProject, " ", todo}
		a.AddTodo(strings.Join(args, ""))
	}

	a.TodoList.Delete(id)
	a.save()
	fmt.Println("Todo expanded.")
}

// HandleNotes is a sub-function that will handle notes on a todo.
func (a *App) HandleNotes(input string) {
	a.Load()
	id := a.getID(input)
	if id == -1 {
		return
	}
	todo := a.TodoList.FindByID(id)
	if todo == nil {
		fmt.Println("No such id.")
		return
	}
	parser := &Parser{}

	if parser.ParseAddNote(todo, input) {
		fmt.Println("Note added.")
	} else if parser.ParseDeleteNote(todo, input) {
		fmt.Println("Note deleted.")
	} else if parser.ParseEditNote(todo, input) {
		fmt.Println("Note edited.")
	} else if parser.ParseShowNote(todo, input) {
		groups := map[string][]*Todo{}
		groups[""] = append(groups[""], todo)
		a.Printer.Print(&GroupedTodos{Groups: groups}, true)
		return
	}
	a.save()
}

// ArchiveCompleted will archive all completed todos.
func (a *App) ArchiveCompleted() {
	a.Load()
	for _, todo := range a.TodoList.Todos() {
		if todo.Completed {
			todo.Archive()
		}
	}
	a.save()
	fmt.Println("All completed todos have been archived.")
}

// ListTodos will list all todos.
func (a *App) ListTodos(input string) {
	a.Load()
	filtered := NewFilter(a.TodoList.Todos()).Filter(input)
	grouped := a.getGroups(input, filtered)

	re, _ := regexp.Compile(`^ln`)
	a.Printer.Print(grouped, re.MatchString(input))
}

// PrioritizeTodo will prioritize a todo.
func (a *App) PrioritizeTodo(input string) {
	a.Load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Prioritize(ids...)
	a.save()
	fmt.Println("Todo prioritized.")
}

// UnprioritizeTodo unprioritizes a todo.
func (a *App) UnprioritizeTodo(input string) {
	a.Load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Unprioritize(ids...)
	a.save()
	fmt.Println("Todo un-prioritized.")
}

// GarbageCollect will delete all archived todos.
func (a *App) GarbageCollect() {
	a.Load()
	a.TodoList.GarbageCollect()
	a.save()
	fmt.Println("Garbage collection complete.")
}

// Sync will sync the todolist with ultralist.io.
func (a *App) Sync(quiet bool) {
	a.Load()

	if a.EventLogger.CurrentSyncedList.Name == "" {
		prompt := promptui.Prompt{
			Label: "Give this list a name",
		}

		result, _ := prompt.Run()
		a.EventLogger.CurrentSyncedList.Name = result
	}

	var synchronizer *Synchronizer
	if quiet {
		synchronizer = NewQuietSynchronizer()
	} else {
		synchronizer = NewSynchronizer()
	}
	synchronizer.Sync(a.TodoList, a.EventLogger.CurrentSyncedList)

	if synchronizer.WasSuccessful() {
		a.EventLogger.ClearEventLogs()
		a.TodoStore.Save(a.TodoList)
	}
}

// CheckAuth is checking the authentication against ultralist.io.
func (a *App) CheckAuth() {
	synchronizer := NewSynchronizer()
	synchronizer.CheckAuth()
}

// AuthWorkflow is creating the authentication workflow.
func (a *App) AuthWorkflow() {
	webapp := &Webapp{}
	backend := NewBackend()

	open.Start(backend.AuthURL())
	fmt.Println("Listening for auth response...")
	webapp.Run()
}

// Load the todolist from the todo store.
func (a *App) Load() error {
	todos, err := a.TodoStore.Load()
	if err != nil {
		return err
	}
	a.TodoList.Load(todos)
	a.EventLogger = NewEventLogger(a.TodoList, a.TodoStore)
	a.EventLogger.LoadSyncedLists()
	return nil
}

// OpenWeb is opening the current list on ultralist.io in your browser.
func (a *App) OpenWeb() {
	a.Load()
	if !a.TodoList.IsSynced {
		fmt.Println("This list isn't synced! Use 'ultralist sync' to synchronize this list with ultralist.io.")
		return
	}

	fmt.Println("Opening this list on your browser...")
	open.Start("https://app.ultralist.io/todolist/" + a.EventLogger.CurrentSyncedList.UUID)
}

// Save the todolist to the store.
func (a *App) save() {
	a.TodoStore.Save(a.TodoList)
	if a.TodoList.IsSynced {
		a.EventLogger.ProcessEvents()

		synchronizer := NewQuietSynchronizer()
		synchronizer.ExecSyncInBackground()
	}
}

func (a *App) getID(input string) int {
	re, _ := regexp.Compile("\\d+")
	if re.MatchString(input) {
		id, _ := strconv.Atoi(re.FindString(input))
		return id
	}

	fmt.Println("Invalid id.")
	return -1
}

func (a *App) getIDs(input string) (ids []int) {
	idGroups := strings.Split(input, ",")
	for _, idGroup := range idGroups {
		if rangedIds, err := a.parseRangedIds(idGroup); len(rangedIds) > 0 || err != nil {
			if err != nil {
				fmt.Printf("Invalid id group: %s.\n", input)
				continue
			}
			ids = append(ids, rangedIds...)
		} else if id := a.getID(idGroup); id != -1 {
			ids = append(ids, id)
		} else {
			fmt.Printf("Invalid id: %s.\n", idGroup)
		}
	}
	return ids
}

func (a *App) parseRangedIds(input string) (ids []int, err error) {
	rangeNumberRE, _ := regexp.Compile("(\\d+)-(\\d+)")
	if matches := rangeNumberRE.FindStringSubmatch(input); len(matches) > 0 {
		lowerID, _ := strconv.Atoi(matches[1])
		upperID, _ := strconv.Atoi(matches[2])
		if lowerID >= upperID {
			return ids, fmt.Errorf("Invalid id group: %s", input)
		}
		for id := lowerID; id <= upperID; id++ {
			ids = append(ids, id)
		}
	}
	return ids, err
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
