package ultralist

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Webapp struct {
	Router *httprouter.Router
	server *http.Server
}

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
	keys, ok := r.URL.Query()["token"]
	if ok == false {
		fmt.Println("Something went wrong.. I did not get a token back.")
		os.Exit(0)
	}

	backend := NewBackend()
	backend.WriteCreds(keys[0])
	fmt.Println("Authorization successful!")
	fmt.Fprintf(writer, "Authorization complete.  Head back to your terminal for next steps.")

	// sleep 1 second before shutting server down, so we can display msg on web.
	go func() {
		time.Sleep(1 * time.Second)
		w.server.Shutdown(nil)
	}()
}
