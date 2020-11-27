package sync

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jinzhu/copier"
	"github.com/twinj/uuid"
	"github.com/ultralist/ultralist/store"
	"github.com/ultralist/ultralist/ultralist"
)

// The different types of events that can occur.
const (
	AddEvent    = "EventAdded"
	UpdateEvent = "EventUpdated"
	DeleteEvent = "EventDeleted"
)

// EventLogger is the main struct of this file.
type EventLogger struct {
	PreviousTodoList  *ultralist.TodoList
	CurrentTodoList   *ultralist.TodoList
	Store             store.Store
	SyncedLists       []*SyncedList
	CurrentSyncedList *SyncedList
	Events            []*EventLog
}

// SyncedList is a representation of a todolist for syncing.
type SyncedList struct {
	Filename string      `json:"filename"`
	UUID     string      `json:"uuid"`
	Name     string      `json:"name"`
	Events   []*EventLog `json:"events"`
}

// EventLog is a log of events that occurred, with the todo data.
type EventLog struct {
	EventType    string          `json:"event_type"`
	ObjectType   string          `json:"object_type"`
	TodoListUUID string          `json:"todo_list_uuid"`
	Object       *ultralist.Todo `json:"object"`
}

// NewEventLogger is creating a new event logger.
func NewEventLogger(todoList *ultralist.TodoList, store store.Store) *EventLogger {
	var previousTodos []*ultralist.Todo

	for _, todo := range todoList.Data {
		var newTodo ultralist.Todo
		copier.Copy(&newTodo, &todo)
		previousTodos = append(previousTodos, &newTodo)
	}
	var previousTodoList = &ultralist.TodoList{Data: previousTodos}

	eventLogger := &EventLogger{
		CurrentTodoList:  todoList,
		PreviousTodoList: previousTodoList,
		Store:            store,
	}

	eventLogger.loadSyncedLists()

	return eventLogger
}

// ProcessEvents processes all events that occurred when ultralist ran and write them to a log file.
func (e *EventLogger) ProcessEvents() {
	e.CreateEventLogs()
	e.WriteEventLogs()
}

// CreateEventLogs makes event logs.
func (e *EventLogger) CreateEventLogs() {
	var eventLogs []*EventLog

	// find events added or updated
	for _, todo := range e.CurrentTodoList.Data {
		previousTodo := e.PreviousTodoList.FindByID(todo.ID)
		if previousTodo != nil {
			if !todo.Equals(previousTodo) {
				eventLogs = append(eventLogs, e.writeTodoEvent(UpdateEvent, todo, e.CurrentSyncedList.UUID))
			}
		} else {
			eventLogs = append(eventLogs, e.writeTodoEvent(AddEvent, todo, e.CurrentSyncedList.UUID))
		}
	}

	// find deleted events
	for _, todo := range e.PreviousTodoList.Data {
		currentTodo := e.CurrentTodoList.FindByID(todo.ID)
		if currentTodo == nil {
			eventLogs = append(eventLogs, e.writeTodoEvent(DeleteEvent, todo, e.CurrentSyncedList.UUID))
		}
	}
	e.Events = eventLogs
}

// WriteEventLogs writes event logs to disk.
func (e *EventLogger) WriteEventLogs() {
	e.CurrentSyncedList.Events = append(e.CurrentSyncedList.Events, e.Events...)
	e.WriteSyncedLists()
}

// ClearEventLogs is clearing the event logs and writes it to disk.
func (e *EventLogger) ClearEventLogs() {
	e.CurrentSyncedList.Events = []*EventLog{}
	e.WriteSyncedLists()
}

// LoadSyncedLists - load all currently synced lists from ~/.config/ultralist/synced_lists.json
// marshal them into SyncedList structs, and store them in []SyncedLists
// if the current list is not in that file, then initialize a new synced list with the info we know about the current todo list.
func (e *EventLogger) loadSyncedLists() {
	if _, err := os.Stat(e.syncedListsFile()); os.IsNotExist(err) {
		e.initializeSyncedListFromCurrentTodoList()
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

	e.initializeSyncedListFromCurrentTodoList()
}

func (e *EventLogger) initializeSyncedListFromCurrentTodoList() {

	listUUID := e.CurrentTodoList.UUID
	if listUUID == "" {
		listUUID = fmt.Sprintf("%s", uuid.NewV4())
		e.CurrentTodoList.UUID = listUUID
	}

	list := &SyncedList{
		Filename: e.Store.GetLocation(),
		Name:     e.CurrentTodoList.Name,
		UUID:     listUUID,
	}
	e.SyncedLists = append(e.SyncedLists, list)
	e.CurrentSyncedList = list
}

// DeleteCurrentSyncedList - delete a synced list from the synced_lists.json file
func (e *EventLogger) DeleteCurrentSyncedList() {
	var syncedListsWithoutDeleted []*SyncedList
	for _, list := range e.SyncedLists {
		if list.UUID != e.CurrentSyncedList.UUID {
			syncedListsWithoutDeleted = append(syncedListsWithoutDeleted, list)
		}
	}
	e.SyncedLists = syncedListsWithoutDeleted
}

// WriteSyncedLists is writing a synced list.
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
	home := e.userHomeDir()
	return fmt.Sprintf("%s/.config/ultralist/", home)
}

func (e *EventLogger) syncedListsFile() string {
	return e.syncedListsConfigDir() + "synced_lists.json"
}

func (e *EventLogger) writeTodoEvent(eventType string, todo *ultralist.Todo, todoListUUID string) *EventLog {
	return &EventLog{
		EventType:    eventType,
		ObjectType:   "TodoItem",
		TodoListUUID: todoListUUID,
		Object:       todo,
	}
}

// UserHomeDir returns the home dir of the current user.
func (e *EventLogger) userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}
