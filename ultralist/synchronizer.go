package ultralist

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Synchronizer is the default struct of this file.
type Synchronizer struct {
	QuietSync bool
	Success   bool
	Backend   *Backend
}

// NewSynchronizer is creating a new synchronizer.
func NewSynchronizer() *Synchronizer {
	return &Synchronizer{QuietSync: false, Success: false, Backend: NewBackend()}
}

func NewQuietSynchronizer() *Synchronizer {
	return &Synchronizer{QuietSync: true, Success: false, Backend: NewBackend()}
}

// NewSynchronizerWithInput is creating a new synchronizer with input.
func NewSynchronizerWithInput(input string) *Synchronizer {
	quietSync := false
	if input == "sync -q" {
		quietSync = true
	}
	return &Synchronizer{QuietSync: quietSync, Success: false, Backend: NewBackend()}
}

// ExecSyncInBackground is starting a new sync process with the ultralist API in the background.
func (s *Synchronizer) ExecSyncInBackground() {
	binary, lookErr := exec.LookPath("ultralist")
	if lookErr != nil {
		panic(lookErr)
	}

	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	process, err := os.StartProcess(binary, []string{binary, "sync", "q"}, &procAttr)

	if err != nil {
		panic(err)
	}

	err = process.Release()
	if err != nil {
		panic(err)
	}
}

// Sync is synchronizing the todos with the ultralist API.
func (s *Synchronizer) Sync(todolist *TodoList, syncedList *SyncedList) {

	if s.Backend.CredsFileExists() == false {
		s.println("Cannot find credentials file.  Please re-authorize!")
		return
	}

	if s.Backend.CanConnect() == false {
		s.println("Cannot connect to api.ultralist.io right now.")
		return
	}

	s.doSync(todolist, syncedList)
}

// CheckAuth is checking the authentication status against the ultralist API.
func (s *Synchronizer) CheckAuth() {
	if s.Backend.CredsFileExists() == false {
		fmt.Println("It looks like you are not authenticated with ultralist.io.")
		return
	}

	if s.Backend.CanConnect() == false {
		fmt.Println("Cannot connect to api.ultralist.io right now.")
		return
	}

	bodyBytes := s.Backend.PerformRequest("GET", "/me", []byte(""))

	var response *UserRequest
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		panic(err)
	}
	if s.Backend.Success {
		s.Success = true
		fmt.Printf("Hello %s! You are successfully authenticated.\n", response.Name)
	}
}

// WasSuccessful is checking if a sync process was successful.
func (s *Synchronizer) WasSuccessful() bool {
	return s.Success
}

// UserRequest is the struct for a user request.
type UserRequest struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// TodolistRequest is the struct for a todolist request.
type TodolistRequest struct {
	UUID                string  `json:"uuid"`
	Name                string  `json:"name"`
	TodoItemsAttributes []*Todo `json:"todo_items_attributes"`
}

// Request is the struct for a request.
type Request struct {
	Events   []*EventLog      `json:"events"`
	From     string           `json:"from"`
	Todolist *TodolistRequest `json:"todolist"`
}

func (s *Synchronizer) doSync(todolist *TodoList, syncedList *SyncedList) {
	data := s.buildRequest(todolist, syncedList)
	path := "/api/v1/todo_lists/event_cache"

	s.Backend.PerformRequest("PUT", path, data)

	// after performing the request, must do GET request to the todo list in question and fill it out
	// the server will have the "correct" list, since it will have assimilated all of the change
	// from various clients.

	path = "/api/v1/todo_lists/" + syncedList.UUID
	bodyBytes := s.Backend.PerformRequest("GET", path, data)

	var response *TodolistRequest
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		panic(err)
	}
	if s.Backend.Success {
		s.Success = true
		todolist.Data = response.TodoItemsAttributes
	}
}

func (s *Synchronizer) buildRequest(todolist *TodoList, syncedList *SyncedList) []byte {
	requestData := &Request{
		From:   "cli",
		Events: syncedList.Events,
		Todolist: &TodolistRequest{
			UUID: syncedList.UUID,
			Name: syncedList.Name,
		},
	}
	data, _ := json.Marshal(requestData)
	return data
}

func (s *Synchronizer) println(text string) {
	if s.QuietSync == false {
		fmt.Println(text)
	}
}
