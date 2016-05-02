package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gammons/todolist/todolist"
)

func main() {
	if len(os.Args) <= 1 {
		usage()
		os.Exit(0)
	}
	input := strings.Join(os.Args[1:], " ")
	routeInput(os.Args[1], input)
}

func usage() {
	fmt.Println("usage")
}

func routeInput(command string, input string) {
	app := todolist.NewApp()
	switch {
	case command == "l" || command == "list":
		app.ListTodos(input)
	case command == "a" || command == "add":
		app.AddTodo(input)
	case command == "d" || command == "del":
		app.DeleteTodo(input)
	case command == "c" || command == "complete":
		app.CompleteTodo(input)
	case command == "uc" || command == "uncomplete":
		app.UncompleteTodo(input)
	case command == "ar" || command == "archive":
		app.ArchiveTodo(input)
	case command == "uar" || command == "unarchive":
		app.UnarchiveTodo(input)
	}
}
