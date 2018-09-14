package main

import (
	"os"
	"strings"

	"github.com/gammons/todolist/todolist"
)

// the current version of todolist
const (
	VERSION = "0.8.1"
)

func main() {
	if len(os.Args) <= 1 {
		todolist.Usage()
		os.Exit(0)
	}
	input := strings.Join(os.Args[1:], " ")
	routeInput(os.Args[1], input)
}

func routeInput(command string, input string) {
	app := todolist.NewApp()
	switch command {
	case "l", "ln", "list", "agenda":
		app.ListTodos(input)
	case "a", "add":
		app.AddTodo(input)
	case "done":
		app.AddDoneTodo(input)
	case "d", "delete":
		app.DeleteTodo(input)
	case "c", "complete":
		app.CompleteTodo(input)
	case "uc", "uncomplete":
		app.UncompleteTodo(input)
	case "ar", "archive":
		app.ArchiveTodo(input)
	case "uar", "unarchive":
		app.UnarchiveTodo(input)
	case "ac":
		app.ArchiveCompleted()
	case "e", "edit":
		app.EditTodo(input)
	case "ex", "expand":
		app.ExpandTodo(input)
	case "an", "n", "dn", "en":
		app.HandleNotes(input)
	case "gc":
		app.GarbageCollect()
	case "p", "prioritize":
		app.PrioritizeTodo(input)
	case "up", "unprioritize":
		app.UnprioritizeTodo(input)
	case "init":
		app.InitializeRepo()
	case "sync":
		app.Sync(input)
	case "auth":
		app.AuthWorkflow()
	}
}
