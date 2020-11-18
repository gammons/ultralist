package ultralist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	apiURL = "https://api.ultralist.io"
)

// Backend is giving you the structure of the ultralist API backend.
type Backend struct {
	Creds   string `json:"creds"`
	Success bool   `json:"-"`
}

// NewBackend is starting a new backend.
func NewBackend() *Backend {
	backend := &Backend{Success: false}

	if backend.CredsFileExists() {
		backend.loadCreds()
	}

	return backend
}

// CreateTodoList will create a todo list on the backend
func (b *Backend) CreateTodoList(todolist *TodoList) {
	type Request struct {
		Todolist *TodoList `json:"todolist"`
	}

	bodyBytes, _ := json.Marshal(&Request{Todolist: todolist})

	b.PerformRequest("POST", "/api/v1/todo_lists", bodyBytes)
}

// PerformRequest is performing a request to the ultralist API backend.
func (b *Backend) PerformRequest(method string, path string, data []byte) []byte {
	url := b.apiURL(path)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(data))
	authHeader := fmt.Sprintf("Bearer %s", b.Creds)

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-Ultralist-CLI-Version", VERSION)

	client := &http.Client{}

	var requestError error
	var response *http.Response

	response, requestError = client.Do(req)
	//defer response.Body.Close()

	if requestError != nil {
		fmt.Println("Error contacting server: ", requestError)
		b.Success = false
		return nil
	}

	if !strings.HasPrefix(response.Status, "2") {
		fmt.Printf("Got a status of %s from server. Aborting.", response.Status)
		os.Exit(0)
		return nil
	}

	b.Success = true

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return bodyBytes
}

// CanConnect is checking if ultralist can connect to the API backend.
func (b *Backend) CanConnect() bool {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{Timeout: timeout}
	if _, err := client.Get(b.apiURL("/hb")); err != nil {
		return false
	}
	return true
}

// CredsFileExists is checking if an credential file is existing.
func (b *Backend) CredsFileExists() bool {
	if _, err := os.Stat(b.credsFilePath()); os.IsNotExist(err) {
		return false
	}
	return true
}

// AuthURL is providing the full auth URL for a API backend.
func (b *Backend) AuthURL() string {
	return b.apiURL("/cli_auth")
}

// WriteCreds is writing the backend credential file.
func (b *Backend) WriteCreds(token string) {
	b.Creds = token
	data, _ := json.Marshal(b)

	if _, err := os.Stat(b.credsFolderPath()); os.IsNotExist(err) {
		if err := os.MkdirAll(b.credsFolderPath(), os.ModePerm); err != nil {
			fmt.Println("Could not create ~/.config/ultralist directory! Permissions issue?")
			os.Exit(1)
		}
	}

	if err := ioutil.WriteFile(b.credsFilePath(), data, 0600); err != nil {
		fmt.Println("Error writing creds file!")
		os.Exit(1)
	}
}

func (b *Backend) apiURL(path string) string {
	envAPIURL := os.Getenv("ULTRALIST_API_URL")

	if envAPIURL != "" {
		return envAPIURL + path
	}
	return apiURL + path
}

func (b *Backend) credsFolderPath() string {
	home := UserHomeDir()
	return fmt.Sprintf("%s/.config/ultralist/", home)
}

func (b *Backend) credsFilePath() string {
	return b.credsFolderPath() + "creds.json"
}

func (b *Backend) loadCreds() string {
	data, err := ioutil.ReadFile(b.credsFilePath())
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, b)
	if err != nil {
		panic(err)
	}
	return string(data)
}
