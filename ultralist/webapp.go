package ultralist

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Webapp is the main struct of this file.
type Webapp struct {
	Router *httprouter.Router
	server *http.Server
}

// Run is starting the ultralist webapp.
func (w *Webapp) Run() {
	w.server = &http.Server{Addr: ":9976"}

	http.HandleFunc("/", w.handleAuthResponse)
	http.HandleFunc("/favicon.ico", w.handleFavicon)

	w.server.ListenAndServe()
}

func (w *Webapp) handleFavicon(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(writer, "")
}

func (w *Webapp) handleAuthResponse(writer http.ResponseWriter, r *http.Request) {
	cliTokens, ok := r.URL.Query()["token"]
	if ok == false {
		fmt.Println("Something went wrong... I did not get a CLI token back.")
		os.Exit(0)
	}

	backend := NewBackend()
	backend.WriteCreds(cliTokens[0])
	fmt.Println("Authorization successful! Next, run `ultralist sync --setup` to sync a list.")

	http.Redirect(writer, r, w.frontendUrl(), http.StatusSeeOther)

	// sleep 1 second before shutting server down, so we can display msg on web.
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

func (w *Webapp) frontendUrl() string {
	envFrontendURL := os.Getenv("ULTRALIST_FRONTEND_URL")

	if envFrontendURL != "" {
		return envFrontendURL + "/todolist?cli_auth_completed=true"
	}

	return "https://app.ultralist.io/todolist?cli_auth_completed=true"
}
