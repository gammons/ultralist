package ultralist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jinzhu/copier"
)

// The different types of events that can occur.
const (
	AddEvent    = "EventAdded"
	UpdateEvent = "EventUpdated"
	DeleteEvent = "EventDeleted"
)

// EventLogger is the main struct of this file.
type EventLogger struct {
	PreviousTodoList  *TodoList
	CurrentTodoList   *TodoList
	Store             Store
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
	Status        string   `json:"status"`
	IsPriority    bool     `json:"isPriority"`
	Notes         []string `json:"notes"`
}

// NewEventLogger is creating a new event logger.
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
			if todo.Equals(previousTodo) == false {
				eventLogs = append(eventLogs, e.writeTodoEvent(UpdateEvent, todo))
			}
		} else {
			eventLogs = append(eventLogs, e.writeTodoEvent(AddEvent, todo))
		}
	}

	// find deleted events
	for _, todo := range e.PreviousTodoList.Data {
		currentTodo := e.CurrentTodoList.FindByID(todo.ID)
		if currentTodo == nil {
			eventLogs = append(eventLogs, e.writeTodoEvent(DeleteEvent, todo))
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

// LoadSyncedLists is loading a synced list.
func (e *EventLogger) LoadSyncedLists() {
	if _, err := os.Stat(e.syncedListsFile()); os.IsNotExist(err) {
		e.initializeSyncedList()
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

	e.initializeSyncedList()
}

func (e *EventLogger) initializeSyncedList() {
	list := &SyncedList{
		Filename: e.Store.GetLocation(),
		UUID:     newUUID(),
	}
	e.SyncedLists = append(e.SyncedLists, list)
	e.CurrentSyncedList = list
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
	home := UserHomeDir()
	return fmt.Sprintf("%s/.config/ultralist/", home)
}

func (e *EventLogger) syncedListsFile() string {
	return e.syncedListsConfigDir() + "synced_lists.json"
}

func (e *EventLogger) writeTodoEvent(eventType string, todo *Todo) *EventLog {
	return &EventLog{
		EventType:     eventType,
		ID:            todo.ID,
		UUID:          todo.UUID,
		Subject:       todo.Subject,
		Projects:      todo.Projects,
		Contexts:      todo.Contexts,
		Due:           todo.Due,
		Completed:     todo.Completed,
		CompletedDate: todo.CompletedDate,
		Archived:      todo.Archived,
		Status:        todo.Status,
		IsPriority:    todo.IsPriority,
		Notes:         todo.Notes,
	}
}
