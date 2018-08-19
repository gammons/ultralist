package todolist

import (
	"net/http"
	"os"
	"time"
)

type Synchronizer struct{}

const (
	ApiUrl = "https://api.ultralist.io"
)

func (s *Synchronizer) Sync(todolist *TodoList) {
	if s.canConnect() == false {
		return
	}
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
