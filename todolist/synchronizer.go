package todolist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"time"
)

type Synchronizer struct {
	Creds     string `json:"creds"`
	QuietSync bool
}

const (
	ApiUrl = "https://api.ultralist.io"
)

func NewSynchronizer(input string) *Synchronizer {
	quietSync := false
	if input == "sync -q" {
		quietSync = true
	}

	return &Synchronizer{QuietSync: quietSync}
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

	s.sync(todolist, syncedList)
}

type RequestTodolist struct {
	UUID                string  `json:"uuid"`
	Name                string  `json:"name"`
	TodoItemsAttributes []*Todo `json:"todo_items_attributes"`
}
type Request struct {
	Events   []*EventLog      `json:"events"`
	Todolist *RequestTodolist `json:"todolist"`
}

func (s *Synchronizer) sync(todolist *TodoList, syncedList *SyncedList) {
	requestData := &Request{
		Events: syncedList.Events,
		Todolist: &RequestTodolist{
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

	if response, requestError = client.Do(req); requestError != nil {
		fmt.Println("Error contacting server: ", requestError)
		os.Exit(0)
	}
	defer response.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	panic(err)
	// }
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
	fmt.Println(string(data))
	return string(data)
}
