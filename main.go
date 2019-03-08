package main

import (
	"github.com/ultralist/ultralist/cmd"
	"github.com/ultralist/ultralist/ultralist"
)

func main() {
	cmd.Execute()
}

func routeInput(command string, input string) {
	app := ultralist.NewApp()
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
	case "check_auth":
		app.CheckAuth()
	case "web":
		app.OpenWeb()
	}
}
