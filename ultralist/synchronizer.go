package ultralist

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"
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
	binary, lookErr := exec.LookPath("todolist")
	if lookErr != nil {
		panic(lookErr)
	}
	cmd := exec.Command(binary, "sync", "-q")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
	if err := cmd.Start(); err != nil {
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

func (s *Synchronizer) WasSuccessful() bool {
	return s.Success
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
