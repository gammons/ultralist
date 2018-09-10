package todolist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"syscall"
	"time"
)

type Synchronizer struct {
	Creds     string `json:"creds"`
	QuietSync bool
	Success   bool
}

const (
	ApiUrl = "https://api.ultralist.io"
)

func NewSynchronizer(input string) *Synchronizer {
	quietSync := false
	if input == "sync -q" {
		quietSync = true
	}

	return &Synchronizer{QuietSync: quietSync, Success: false}
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
	if s.credsFileExists() == false {
		s.println("Cannot find credentials file.  Please re-authorize!")
		return
	}

	s.loadCreds()

	if s.canConnect() == false {
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
	bodyBytes := s.performSyncRequest(todolist, syncedList)

	// assign the local todolist data to the values that came back from the server.
	// the server will have the "correct" list, since it will have assimilated all of the change
	// from various clients.
	var response *TodolistRequest
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		panic(err)
	}
	todolist.Data = response.TodoItemsAttributes
}

func (s *Synchronizer) performSyncRequest(todolist *TodoList, syncedList *SyncedList) []byte {
	requestData := &Request{
		Events: syncedList.Events,
		Todolist: &TodolistRequest{
			UUID:                syncedList.UUID,
			Name:                "test todolist",
			TodoItemsAttributes: todolist.Data,
		},
	}

	url := s.apiUrl(fmt.Sprintf("/api/v1/todo_lists/%s", syncedList.UUID))
	data, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	authHeader := fmt.Sprintf("Bearer %s", s.Creds)

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	var requestError error
	var response *http.Response

	response, requestError = client.Do(req)
	defer response.Body.Close()

	if requestError != nil {
		fmt.Println("Error contacting server: ", requestError)
		os.Exit(0)
		return nil
	}

	s.Success = true

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return bodyBytes
}

func (s *Synchronizer) println(text string) {
	if s.QuietSync == false {
		fmt.Println(text)
	}
}

func (s *Synchronizer) credsFileExists() bool {
	if _, err := os.Stat(s.credsFilePath()); os.IsNotExist(err) {
		return false
	}
	return true
}

func (s *Synchronizer) canConnect() bool {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{Timeout: timeout}
	if _, err := client.Get(s.apiUrl("/api/v1/hb")); err != nil {
		return false
	}
	return true
}

func (s *Synchronizer) apiUrl(path string) string {
	apiUrl := os.Getenv("ULTRALIST_API_URL")

	if apiUrl == "" {
		apiUrl = ApiUrl
	}
	return apiUrl + path
}

func (s *Synchronizer) credsFilePath() string {
	usr, _ := user.Current()
	return fmt.Sprintf("%s/.config/ultralist/creds.json", usr.HomeDir)
}

func (s *Synchronizer) loadCreds() string {
	data, err := ioutil.ReadFile(s.credsFilePath())
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, s)
	if err != nil {
		panic(err)
	}
	return string(data)
}
