package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gammons/todolist/todolist"
	"github.com/skratchdot/open-golang/open"
)

const (
	// VERSION specifies the program version
	VERSION = "0.5.3"
)

func main() {
	if len(os.Args) <= 1 {
		usage()
		os.Exit(0)
	}
	input := strings.Join(os.Args[1:], " ")
	err := routeInput(os.Args[1], input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func routeInput(command string, input string) error {
	app := todolist.NewApp()
	switch command {
	case "l", "list", "agenda":
		app.ListTodos(input)
	case "a", "add":
		app.AddTodo(input)
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
	case "e", "edit", "ed", "edit-date":
		app.EditTodoDue(input)
	case "es", "edit-subject":
		app.EditTodoSubject(input)
	case "ex", "expand":
		app.ExpandTodo(input)
	case "gc":
		app.GarbageCollect()
	case "p", "prioritize":
		app.PrioritizeTodo(input)
	case "up", "unprioritize":
		app.UnprioritizeTodo(input)
	case "init":
		app.InitializeRepo()
	case "web":
		if err := app.Load(); err != nil {
			return err
		} else {
			web := todolist.NewWebapp()
			fmt.Println("Now serving todolist web.\nHead to http://localhost:7890 to see your todo list!")
			open.Start("http://localhost:7890")
			web.Run()
		}
	}
	return nil
}
