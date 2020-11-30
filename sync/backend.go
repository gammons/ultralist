package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ultralist/ultralist/ultralist"
)

const (
	apiURL = "https://api.ultralist.io"
)

// Backend is giving you the structure of the ultralist API backend.
type Backend struct {
	Creds string `json:"creds"`
}

// NewBackend is starting a new backend.
func NewBackend() *Backend {
	backend := &Backend{}

	if backend.CredsFileExists() {
		backend.loadCreds()
	}

	return backend
}

// CreateTodoList will create a todo list on the backend
func (b *Backend) CreateTodoList(todolist *ultralist.TodoList) {
	type Request struct {
		Todolist *ultralist.TodoList `json:"todolist"`
	}

	bodyBytes, _ := json.Marshal(&Request{Todolist: todolist})

	b.PerformRequest("POST", "/api/v1/todo_lists", bodyBytes)
}

// PerformRequest is performing a request to the ultralist API backend.
func (b *Backend) PerformRequest(method string, path string, data []byte) ([]byte, error) {
	url := b.apiURL(path)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(data))
	authHeader := fmt.Sprintf("Bearer %s", b.Creds)

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-Ultralist-CLI-Version", ultralist.Version)

	client := &http.Client{}

	var requestError error
	var response *http.Response

	response, requestError = client.Do(req)
	//defer response.Body.Close()

	if requestError != nil {
		return nil, requestError
	}

	if !strings.HasPrefix(response.Status, "2") {
		return nil, fmt.Errorf(fmt.Sprintf("Got a status of %s from server", response.Status))
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

// CanConnect checks to see if it can connect to api.ultralist.io, with a 2 second timeout.
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

// WriteCreds writes the backend credential file.
func (b *Backend) WriteCreds(token string) error {
	b.Creds = token
	data, _ := json.Marshal(b)

	if _, err := os.Stat(b.credsFolderPath()); os.IsNotExist(err) {
		if err := os.MkdirAll(b.credsFolderPath(), os.ModePerm); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(b.credsFilePath(), data, 0600); err != nil {
		return err
	}
	return nil
}

func (b *Backend) apiURL(path string) string {
	envAPIURL := os.Getenv("ULTRALIST_API_URL")

	if envAPIURL != "" {
		return envAPIURL + path
	}
	return apiURL + path
}

func (b *Backend) credsFolderPath() string {
	home, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.config/ultralist/", home)
}

func (b *Backend) credsFilePath() string {
	return b.credsFolderPath() + "creds.json"
}

func (b *Backend) loadCreds() (string, error) {
	data, err := ioutil.ReadFile(b.credsFilePath())
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(data, b)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
