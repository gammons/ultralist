package sync

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/ultralist/ultralist/ultralist"
)

// Synchronizer is the default struct of this file.
type Synchronizer struct {
	QuietSync bool
	Backend   *Backend
}

// NewSynchronizer is creating a new synchronizer.
func NewSynchronizer() *Synchronizer {
	return &Synchronizer{QuietSync: false, Backend: NewBackend()}
}

func NewQuietSynchronizer() *Synchronizer {
	return &Synchronizer{QuietSync: true, Backend: NewBackend()}
}

// ExecSyncInBackground starts a new sync process with the ultralist API in the background.
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
func (s *Synchronizer) Sync(todolist *ultralist.TodoList, syncedList *SyncedList) error {

	if s.Backend.CredsFileExists() == false {
		return errors.New("cannot find credentials file.  Please re-authorize")
	}

	if s.Backend.CanConnect() == false {
		return errors.New("cannot connect to api.ultralist.io right now")
	}

	if err := s.doSync(todolist, syncedList); err != nil {
		return err
	}

	return nil
}

// CheckAuth is checking the authentication status against the ultralist API.
func (s *Synchronizer) CheckAuth() (string, error) {
	if s.Backend.CredsFileExists() == false {
		return "", errors.New("It looks like you are not authenticated with ultralist.io")
	}

	if s.Backend.CanConnect() == false {
		return "", errors.New("cannot connect to api.ultralist.io right now")
	}

	bodyBytes, err := s.Backend.PerformRequest("GET", "/me", []byte(""))
	if err != nil {
		return "", err
	}

	var response *UserRequest
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", err
	}

	return response.Name, nil
}

// UserRequest is the struct for a user request.
type UserRequest struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// TodolistRequest is the struct for a todolist request.
type TodolistRequest struct {
	UUID                string            `json:"uuid"`
	Name                string            `json:"name"`
	TodoItemsAttributes []*ultralist.Todo `json:"todo_items_attributes"`
}

// Request is the struct for a request.
type Request struct {
	Events   []*EventLog      `json:"events"`
	From     string           `json:"from"`
	Todolist *TodolistRequest `json:"todolist"`
}

func (s *Synchronizer) doSync(todolist *ultralist.TodoList, syncedList *SyncedList) error {
	data := s.buildRequest(todolist, syncedList)
	path := "/api/v1/todo_lists/event_cache"

	s.Backend.PerformRequest("PUT", path, data)

	// after performing the request, must do GET request to the todo list in question and fill it out
	// the server will have the "correct" list, since it will have assimilated all of the change
	// from various clients.

	path = "/api/v1/todo_lists/" + syncedList.UUID
	bodyBytes, err := s.Backend.PerformRequest("GET", path, []byte{})
	if err != nil {
		return err
	}

	var response *TodolistRequest
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return err
	}
	todolist.Data = response.TodoItemsAttributes

	return nil
}

func (s *Synchronizer) buildRequest(todolist *ultralist.TodoList, syncedList *SyncedList) []byte {
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
