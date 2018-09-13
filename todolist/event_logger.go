package todolist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/jinzhu/copier"
)

// The different types of events that can occur
const (
	AddEvent    = "EventAdded"
	UpdateEvent = "EventUpdated"
	DeleteEvent = "EventDeleted"
)

// EventLogger is the main struct of this file
type EventLogger struct {
	PreviousTodoList  *TodoList
	CurrentTodoList   *TodoList
	Store             Store
	SyncedLists       []*SyncedList
	CurrentSyncedList *SyncedList
	Events            []*EventLog
}

// SyncedList is a representation of a todolist for syncing
type SyncedList struct {
	Filename string      `json:"filename"`
	UUID     string      `json:"uuid"`
	Name     string      `json:"name"`
	Events   []*EventLog `json:"events"`
}

// EventLog is a log of events that occurred, with the Todo data.
type EventLog struct {
	EventType     string   `json:"eventType"`
	ID            int      `json:"id"`
	UUID          string   `json:"uuid"`
	Subject       string   `json:"subject"`
	Projects      []string `json:"projects"`
	Contexts      []string `json:"contexts"`
	Due           string   `json:"due"`
	Completed     bool     `json:"completed"`
	CompletedDate string   `json:"completedDate"`
	Archived      bool     `json:"archived"`
	IsPriority    bool     `json:"isPriority"`
	Notes         []string `json:"notes"`
}

// NewEventLogger : Create a new event logger
func NewEventLogger(todoList *TodoList, store Store) *EventLogger {
	var previousTodos []*Todo

	for _, todo := range todoList.Data {
		var newTodo Todo
		copier.Copy(&newTodo, &todo)
		previousTodos = append(previousTodos, &newTodo)
	}
	var previousTodoList = &TodoList{Data: previousTodos}
	return &EventLogger{
		CurrentTodoList:  todoList,
		PreviousTodoList: previousTodoList,
		Store:            store,
	}
}

// ProcessEvents : process all events that occurred when todolist ran, and write them to a log file.
func (e *EventLogger) ProcessEvents() {
	e.CreateEventLogs()
	e.WriteEventLogs()
}

// CreateEventLogs : makes event logs.
func (e *EventLogger) CreateEventLogs() {
	var eventLogs []*EventLog

	// find events added or updated
	for _, todo := range e.CurrentTodoList.Data {
		previousTodo := e.PreviousTodoList.FindById(todo.Id)
		if previousTodo != nil {
			if todo.Equals(previousTodo) == false {
				eventLogs = append(eventLogs, e.writeTodoEvent(UpdateEvent, todo))
			}
		} else {
			eventLogs = append(eventLogs, e.writeTodoEvent(AddEvent, todo))
		}
	}

	// find deleted events
	for _, todo := range e.PreviousTodoList.Data {
		currentTodo := e.CurrentTodoList.FindById(todo.Id)
		if currentTodo == nil {
			eventLogs = append(eventLogs, e.writeTodoEvent(DeleteEvent, todo))
		}
	}
	e.Events = eventLogs
}

// WriteEventLogs : Writes event logs to disk
func (e *EventLogger) WriteEventLogs() {
	e.LoadSyncedLists()
	e.CurrentSyncedList.Events = append(e.CurrentSyncedList.Events, e.Events...)
	e.WriteSyncedLists()
}

func (e *EventLogger) ClearEventLogs() {
	e.LoadSyncedLists()
	e.CurrentSyncedList.Events = []*EventLog{}
	e.WriteSyncedLists()
}

func (e *EventLogger) LoadSyncedLists() {
	if _, err := os.Stat(e.syncedListsFile()); os.IsNotExist(err) {
		list := &SyncedList{
			Filename: e.Store.GetLocation(),
			UUID:     newUUID(),
		}
		e.SyncedLists = []*SyncedList{list}
		e.CurrentSyncedList = list
		return
	}

	data, _ := ioutil.ReadFile(e.syncedListsFile())
	err := json.Unmarshal(data, &e.SyncedLists)
	if err != nil {
		panic(err)
	}

	for _, list := range e.SyncedLists {
		if list.Filename == e.Store.GetLocation() {
			e.CurrentSyncedList = list
			e.CurrentTodoList.IsSynced = true
			return
		}
	}
	e.CurrentSyncedList = &SyncedList{
		Filename: e.Store.GetLocation(),
		UUID:     newUUID(),
	}
	e.SyncedLists = append(e.SyncedLists, e.CurrentSyncedList)
}

func (e *EventLogger) WriteSyncedLists() {
	data, _ := json.Marshal(e.SyncedLists)
	if _, err := os.Stat(e.syncedListsConfigDir()); os.IsNotExist(err) {
		os.MkdirAll(e.syncedListsConfigDir(), os.ModePerm)
		if _, cerr := os.Create(e.syncedListsConfigDir()); cerr != nil {
			panic(cerr)
		}
	}

	if err := ioutil.WriteFile(e.syncedListsFile(), data, 0644); err != nil {
		panic(err)
	}
}

func (e *EventLogger) syncedListsConfigDir() string {
	usr, _ := user.Current()
	return fmt.Sprintf("%s/.config/ultralist/", usr.HomeDir)
}

func (e *EventLogger) syncedListsFile() string {
	return e.syncedListsConfigDir() + "synced_lists.json"
}

func (e *EventLogger) writeTodoEvent(eventType string, todo *Todo) *EventLog {
	return &EventLog{
		EventType:     eventType,
		ID:            todo.Id,
		UUID:          todo.UUID,
		Subject:       todo.Subject,
		Projects:      todo.Projects,
		Contexts:      todo.Contexts,
		Due:           todo.Due,
		Completed:     todo.Completed,
		CompletedDate: todo.CompletedDate,
		Archived:      todo.Archived,
		IsPriority:    todo.IsPriority,
		Notes:         todo.Notes,
	}
}
