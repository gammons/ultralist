package todolist

import (
	"fmt"
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

func (s *Synchronizer) Sync(syncedList *SyncedList) {
	if s.credsFileExists() == false {
		s.println("Cannot find credentials file.  Please re-authorize!")
		return
	}

	if s.canConnect() == false {
		s.println("Cannot connect to api.ultralist.io right now.")
		return
	}
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
