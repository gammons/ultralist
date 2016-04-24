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
	routeInput(os.Args[1])
}

func usage() {
	fmt.Println("usage")
}

func routeInput(command string) {
	app := todolist.NewApp()
	input := strings.Join(os.Args[1:], " ")
	switch {
	case command == "l" || command == "list":
		app.ListTodos(input)
	}

}
