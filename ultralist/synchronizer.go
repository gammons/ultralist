package ultralist

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Synchronizer struct {
	QuietSync bool
	Success   bool
	Backend   *Backend
}

const (
	ApiUrl = "https://api.ultralist.io"
)

func NewSynchronizer() *Synchronizer {
	return &Synchronizer{QuietSync: true, Success: false, Backend: NewBackend()}
}

func NewSynchronizerWithInput(input string) *Synchronizer {
	quietSync := false
	if input == "sync -q" {
		quietSync = true
	}
	return &Synchronizer{QuietSync: quietSync, Success: false, Backend: NewBackend()}
}

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

func (s *Synchronizer) WasSuccessful() bool {
	return s.Success
}

type UserRequest struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type TodolistRequest struct {
	UUID                string  `json:"uuid"`
	Name                string  `json:"name"`
	TodoItemsAttributes []*Todo `json:"todo_items_attributes"`
}

type Request struct {
	Events   []*EventLog      `json:"events"`
	Todolist *TodolistRequest `json:"todolist"`
}

func (s *Synchronizer) doSync(todolist *TodoList, syncedList *SyncedList) {
	data := s.buildRequest(todolist, syncedList)
	path := fmt.Sprintf("/api/v1/todo_lists/%s", syncedList.UUID)

	bodyBytes := s.Backend.PerformRequest("PUT", path, data)

	// assign the local todolist data to the values that came back from the server.
	// the server will have the "correct" list, since it will have assimilated all of the change
	// from various clients.
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
		Events: syncedList.Events,
		Todolist: &TodolistRequest{
			UUID:                syncedList.UUID,
			Name:                syncedList.Name,
			TodoItemsAttributes: todolist.Data,
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
