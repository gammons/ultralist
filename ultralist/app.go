package ultralist

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/skratchdot/open-golang/open"
)

// Current version of ultralist.
const (
	VERSION     string = "1.7.0"
	DATE_FORMAT string = "2006-01-02"
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
		TodoList:  &TodoList{},
		Printer:   NewScreenPrinter(true),
		TodoStore: NewFileStore(),
	}
	return app
}

// NewAppWithPrintOptions creates a new app with options for printing on screen.
func NewAppWithPrintOptions(unicodeSupport bool, colorSupport bool) *App {
	var printer Printer
	if colorSupport {
		printer = NewScreenPrinter(unicodeSupport)
	} else {
		printer = NewSimpleScreenPrinter(unicodeSupport)
	}

	app := &App{
		TodoList:  &TodoList{},
		Printer:   printer,
		TodoStore: NewFileStore(),
	}
	return app
}

// InitializeRepo is initializing ultralist repo.
func (a *App) InitializeRepo() {
	a.TodoStore.Initialize()
	fmt.Println("Repo initialized.")
}

// AddTodo adds a new todo to the current list
func (a *App) AddTodo(input string) {
	a.load()
	parser := &InputParser{}

	filter, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("I need more information. Try something like 'todo a chat with @bob due tom'")
		return
	}

	todoItem, err := CreateTodo(filter)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	a.TodoList.Add(todoItem)
	a.save()
	fmt.Printf("Todo %d added.\n", todoItem.ID)
}

// DeleteTodo deletes a todo.
func (a *App) DeleteTodo(input string) {
	a.load()
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
	a.load()
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
	a.load()
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
	a.load()
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
	a.load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Unarchive(ids...)
	a.save()
	fmt.Println("Todo unarchived.")
}

// EditTodo edits a todo with the given input.
func (a *App) EditTodo(todoID int, input string) {
	a.load()
	todo := a.TodoList.FindByID(todoID)
	if todo == nil {
		fmt.Println("No todo with that id.")
		return
	}

	parser := &InputParser{}
	filter, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = EditTodo(todo, a.TodoList, filter); err != nil {
		fmt.Println(err.Error())
		return
	}

	a.save()
	fmt.Println("Todo updated.")
}

// AddNote adds a note to a todo.
func (a *App) AddNote(todoID int, note string) {
	a.load()

	todo := a.TodoList.FindByID(todoID)
	if todo == nil {
		fmt.Println("No todo with that id.")
		return
	}
	todo.Notes = append(todo.Notes, note)

	fmt.Println("Note added.")
	a.save()
}

// EditNote edits a todo's note.
func (a *App) EditNote(todoID int, noteID int, note string) {
	a.load()

	todo := a.TodoList.FindByID(todoID)
	if todo == nil {
		fmt.Println("No todo with that id.")
		return
	}

	if noteID >= len(todo.Notes) {
		fmt.Println("No note could be found with that ID.")
		return
	}

	todo.Notes[noteID] = note

	fmt.Println("Note edited.")
	a.save()
}

// DeleteNote deletes a note from a todo.
func (a *App) DeleteNote(todoID int, noteID int) {
	a.load()

	todo := a.TodoList.FindByID(todoID)
	if todo == nil {
		fmt.Println("No todo with that id.")
		return
	}

	if noteID >= len(todo.Notes) {
		fmt.Println("No note could be found with that ID.")
		return
	}

	todo.Notes = append(todo.Notes[:noteID], todo.Notes[noteID+1:]...)

	fmt.Println("Note deleted.")
	a.save()
}

// ArchiveCompleted will archive all completed todos.
func (a *App) ArchiveCompleted() {
	a.load()
	for _, todo := range a.TodoList.Todos() {
		if todo.Completed {
			todo.Archive()
		}
	}
	a.save()
	fmt.Println("All completed todos have been archived.")
}

// ListTodos will list all todos.
func (a *App) ListTodos(input string, showNotes bool, showStatus bool) {
	a.load()

	parser := &InputParser{}

	filter, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	todoFilter := &TodoFilter{Todos: a.TodoList.Todos(), Filter: filter}
	grouped := a.getGroups(input, todoFilter.ApplyFilter())

	a.Printer.Print(grouped, showNotes, showStatus)
}

// PrioritizeTodo will prioritize a todo.
func (a *App) PrioritizeTodo(input string) {
	a.load()
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
	a.load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}
	a.TodoList.Unprioritize(ids...)
	a.save()
	fmt.Println("Todo un-prioritized.")
}

// StartTodo will start a todo.
func (a *App) SetTodoStatus(input string) {
	a.load()
	ids := a.getIDs(input)
	if len(ids) == 0 {
		return
	}

	splitted := strings.Split(input, " ")

	a.TodoList.SetStatus(splitted[len(splitted)-1], ids...)
	a.save()
	fmt.Println("Todo status updated.")
}

// GarbageCollect will delete all archived todos.
func (a *App) GarbageCollect() {
	a.load()
	a.TodoList.GarbageCollect()
	a.save()
	fmt.Println("Garbage collection complete.")
}

// Sync will sync the todolist with ultralist.io.
func (a *App) Sync(quiet bool) {
	backend := NewBackend()
	if !backend.CredsFileExists() {
		fmt.Println("You're not authenticated with ultralist.io yet.  Please run `ultralist auth` first.")
		return
	}

	a.load()
	if !a.TodoList.IsSynced {
		fmt.Println("This list isn't currently syncing with ultralist.io.  Please run `ultralist sync --setup` to set up syncing.")
		return
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
		a.TodoStore.Save(a.TodoList.Data)
	}
}

// SetupSync sets up a todolist to sync with ultralist.io.
func (a *App) SetupSync() {
	backend := NewBackend()
	if !backend.CredsFileExists() {
		fmt.Println("You're not authenticated with ultralist.io yet.  Please run `ultralist auth` first.")
		return
	}

	if a.TodoStore.LocalTodosFileExists() {
		a.load()

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
			a.EventLogger.CurrentSyncedList.Name = result
			a.TodoList.Name = result
			a.TodoList.IsSynced = true
			a.EventLogger.WriteSyncedLists()

			// right here, I need to run a request to create a todo list via the API.
			backend.CreateTodoList(a.TodoList)

			return
		}
	}

	// at this point, it is known that a local todos file does not exist.
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
	prompt2 := promptui.Select{
		Label: "Choose a list to import",
		Items: todolistNames,
	}

	idx, _, _ := prompt2.Run()
	a.TodoList = &response.Todolists[idx]
	a.TodoStore.Initialize()
	a.TodoStore.Save(a.TodoList.Data)

	a.EventLogger = NewEventLogger(a.TodoList, a.TodoStore)
	a.EventLogger.WriteSyncedLists()
}

// Unsync stops a list from syncing with Ultralist.io.
func (a *App) Unsync() {
	backend := NewBackend()
	if !backend.CredsFileExists() {
		fmt.Println("You're not authenticated with ultralist.io yet.  Please run `ultralist auth` first.")
		return
	}

	a.load()

	if !a.TodoList.IsSynced {
		fmt.Println("This list isn't currently syncing with ultralist.io.")
		return
	}

	a.EventLogger.DeleteCurrentSyncedList()
	a.EventLogger.WriteSyncedLists()
	fmt.Println("This list will no longer sync with ultralist.io.  To set up syncing again, run `ultralist sync --setup`.")
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
	fmt.Println("Head to your browser to complete authorization steps.")
	fmt.Println("Listening for response...")
	webapp.Run()
}

// load the todolist from the todo store.
func (a *App) load() error {
	todos, err := a.TodoStore.Load()
	if err != nil {
		return err
	}
	a.TodoList.Load(todos)
	a.EventLogger = NewEventLogger(a.TodoList, a.TodoStore)
	return nil
}

// OpenWeb is opening the current list on ultralist.io in your browser.
func (a *App) OpenWeb() {
	a.load()
	if !a.TodoList.IsSynced {
		fmt.Println("This list isn't synced! Use 'ultralist sync' to synchronize this list with ultralist.io.")
		return
	}

	fmt.Println("Opening this list on your browser...")
	open.Start("https://app.ultralist.io/todolist/" + a.EventLogger.CurrentSyncedList.UUID)
}

// Save the todolist to the store.
func (a *App) save() {
	a.TodoStore.Save(a.TodoList.Data)
	if a.TodoList.IsSynced {
		a.EventLogger.ProcessEvents()

		synchronizer := NewQuietSynchronizer()
		synchronizer.ExecSyncInBackground()
	}
}

func (a *App) getID(input string) (int, error) {
	splitted := strings.Split(input, " ")
	id, err := strconv.Atoi(splitted[0])
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Invalid id: '%s'", splitted[0]))
	}
	return id, nil
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
		} else if id, err := a.getID(idGroup); err == nil {
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
	contextRegex, _ := regexp.Compile(`group:c.*$`)
	projectRegex, _ := regexp.Compile(`group:p.*$`)
	statusRegex, _ := regexp.Compile(`group:s.*$`)

	var grouped *GroupedTodos

	if contextRegex.MatchString(input) {
		grouped = grouper.GroupByContext(todos)
	} else if projectRegex.MatchString(input) {
		grouped = grouper.GroupByProject(todos)
	} else if statusRegex.MatchString(input) {
		grouped = grouper.GroupByStatus(todos)
	} else {
		grouped = grouper.GroupByNothing(todos)
	}
	return grouped
}
