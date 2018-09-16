package todolist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"
)

type Backend struct {
	Creds   string `json:"creds"`
	Success bool   `json:"-"`
}

func NewBackend() *Backend {
	backend := &Backend{Success: false}

	if backend.CredsFileExists() {
		backend.loadCreds()
	}

	return backend
}

func (b *Backend) PerformRequest(method string, path string, data []byte) []byte {
	url := b.apiUrl(path)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(data))
	authHeader := fmt.Sprintf("Bearer %s", b.Creds)

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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

func (b *Backend) CanConnect() bool {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{Timeout: timeout}
	if _, err := client.Get(b.apiUrl("/api/v1/hb")); err != nil {
		return false
	}
	return true
}

func (b *Backend) CredsFileExists() bool {
	if _, err := os.Stat(b.credsFilePath()); os.IsNotExist(err) {
		return false
	}
	return true
}

func (b *Backend) AuthUrl() string {
	return b.apiUrl("/cli_auth")
}

func (b *Backend) WriteCreds(token string) {
	b.Creds = token
	data, _ := json.Marshal(b)

	if err := ioutil.WriteFile(b.credsFilePath(), data, 0600); err != nil {
		fmt.Println("Error writing creds file!")
		panic(err)
	}
}

func (b *Backend) apiUrl(path string) string {
	apiUrl := os.Getenv("ULTRALIST_API_URL")

	if apiUrl == "" {
		apiUrl = ApiUrl
	}
	return apiUrl + path
}

func (b *Backend) credsFilePath() string {
	usr, _ := user.Current()
	return fmt.Sprintf("%s/.config/ultralist/creds.json", usr.HomeDir)
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
