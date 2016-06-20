package todolist

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	VERSION = "0.2.0"
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
	<!DOCTYPE html>
	<html lang="en">
	<head>
	    <style>
	      body {
		font-family: 'Roboto', sans-serif;
		margin: 0px;
	      }
	      .red {
		color: #D50000 !important;
	      }
	      .blue {
		color: #2196F3 !important;
	      }
	    </style>
		<meta charset="utf-8">
		<link href="https://fonts.googleapis.com/css?family=Roboto:400,300,500" rel="stylesheet">
		<meta name="viewport" content="width=device-width,initial-scale=1">
		<meta name="mobile-web-app-capable" content="yes">
		<title>Todolist</title>
		<style>body{font-family:Roboto,sans-serif;margin:0}</style>
	</head>
	<body>
	<div id="app"></div>
	<script type="text/javascript" src="` + urlFor("common.js") + `"></script>
	<script type="text/javascript" src="` + urlFor("main.js") + `"></script>
	</body>
	</html>
	`
	fmt.Fprintf(w, template)
}

func urlFor(file string) string {
	return "http://d26ykxqbemi1kc.cloudfront.net/" + VERSION + "/" + file
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
