package todolist

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Webapp struct {
	Router *httprouter.Router
}

func NewWebapp() *Webapp {
	return &Webapp{Router: setupRoutes()}
}

func (w *Webapp) Run() {
	log.Fatal(http.ListenAndServe(":7890", w.Router))
}

func setupRoutes() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Scaffold)
	router.OPTIONS("/todos", TodoOptions)
	router.GET("/todos", GetTodos)
	router.POST("/todos", SaveTodos)
	return router
}

func Scaffold(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	template := `
	<html>
	<h1>you are here buddy</h1>
	</html>
	`
	fmt.Fprintf(w, template)
}

func GetTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	app := NewApp()
	app.Load()
	json, _ := json.Marshal(app.TodoList.Data)
	fmt.Fprintf(w, string(json))
}
func TodoOptions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "")
}

func SaveTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	decoder := json.NewDecoder(r.Body)
	var todos []*Todo
	err := decoder.Decode(&todos)
	if err != nil {
		log.Fatal("encountered an error parsing json, ", err)
	}
	app := NewApp()
	app.TodoStore.Save(todos)
}
