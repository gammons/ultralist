package cli

import (
	"fmt"

	"github.com/ultralist/ultralist/ultralist"
)

// App is a representation of the ultralist app that is invoked in CLI mode.
// it will output to stdout.
type App struct {
	UltralistApp *ultralist.App
}

// InitializeRepo will initialize a new .todos.json repo and then tell the user.
func (a *App) InitializeRepo() {
	a.UltralistApp.InitializeRepo()
	fmt.Println("Repo initialized.")
}

// AddTodo adds a todo to the current todolist via the CLI
// Takes an `input`, runs it through Ultralist's InputParser, and then adds the todo to the list
func (cli *App) AddTodo(input string) {
	parser := &InputParser{}

	filter, err := parser.Parse(input)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("I need more information. Try something like 'todo a chat with @bob due tom'")
		return
	}

	todoItem, err := ultralist.CreateTodo(filter)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err = cli.UltralistApp.Load(); err != nil {
		fmt.Println("I had trouble loading the .todos.json file.")
		fmt.Println(err.Error())
		return
	}

	cli.UltralistApp.AddTodo(todoItem)
	cli.UltralistApp.Save()

	fmt.Printf("Todo %d added.\n", todoItem.ID)
}
