package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gammons/todolist/todolist"
	"github.com/skratchdot/open-golang/open"
)

const (
	VERSION = "0.3.1"
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
	blue := color.New(color.FgBlue)
	cyan := color.New(color.FgCyan)
	yellow := color.New(color.FgYellow)

	blueBold := blue.Add(color.Bold)

	fmt.Printf("todo v%s, a simple, command line based, GTD-style todo manager\n", VERSION)

	blueBold.Println("\nAdding todos")
	fmt.Println("  the 'a' command adds todos.")
	fmt.Println("  You can also optionally specify a due date.")
	fmt.Println("  Specify a due date by putting 'due <date>' at the end, where <date> is in (tod|today|tom|tomorrow|mon|tue|wed|thu|fri|sat|sun)")
	fmt.Println("\n  Examples for adding a todo:")
	yellow.Println("\ttodo a Meeting with @bob about +importantPrject due today")
	yellow.Println("\ttodo a +work +verify did @john fix the build\\?")

	blueBold.Println("\nListing todos")
	fmt.Println("  When listing todos, you can filter and group the output.\n")

	fmt.Println("  todo l due (tod|today|tom|tomorrow|overdue|this week|next week|mon|tue|wed|thu|fri|sat|sun|none)")
	fmt.Println("  todo l overdue")

	cyan.Println("  Filtering by date:\n")
	yellow.Println("\ttodo l due tod")
	fmt.Println("\tlists all todos due today\n")
	yellow.Println("\ttodo l due tom")
	fmt.Println("\tlists all todos due tomorrow\n")
	yellow.Println("\ttodo l due mon")
	fmt.Println("\tlists all todos due monday\n")
	yellow.Println("\ttodo l overdue")
	fmt.Println("\tlists all todos where the due date is in the past\n")
	yellow.Println("\ttodo agenda")
	fmt.Println("\tlists all todos where the due date is today or in the past\n")

	cyan.Println("  Grouping:")
	fmt.Println("  You can group todos by context or project.")
	yellow.Println("\ttodo l by c")
	fmt.Println("\tlists all todos grouped by context\n")
	yellow.Println("\ttodo l by p")
	fmt.Println("\tlists all todos grouped by project\n")

	cyan.Println("  Grouping and filtering:")
	fmt.Println("  Of course, you can combine grouping and filtering to get a nice formatted list.\n")
	yellow.Println("\ttodo l due today by c")
	fmt.Println("\tlists all todos due today grouped by context\n")
	yellow.Println("\ttodo l +project due this week by c")
	fmt.Println("\tlists all todos due today for +project, grouped by context\n")
	yellow.Println("\ttodo l @frank due tom by p")
	fmt.Println("\tlists all todos due tomorrow concerining @frank for +project, grouped by project\n")

	blueBold.Println("\nCompleting and uncompleting ")
	fmt.Println("Complete and Uncomplete a todo by its Id:\n")
	yellow.Println("\ttodo c 33")
	fmt.Println("\tCompletes a todo with id 33\n")
	yellow.Println("\ttodo uc 33")
	fmt.Println("\tUncompletes a todo with id 33\n")

	blueBold.Println("\nArchiving")
	fmt.Println("You can archive todos once they are done, or if you might come back to them.")
	fmt.Println("By default, todo will only show unarchived todos.\n")
	yellow.Println("\ttodo ar 33")
	fmt.Println("\tArchives a todo with id 33\n")
	yellow.Println("\ttodo ac")
	fmt.Println("\tArchives all completed todos\n")
	yellow.Println("\ttodo l archived")
	fmt.Println("\tlist all archived todos\n")

	blueBold.Println("\nEditing due dates")
	yellow.Println("\ttodo e 33 due mon")
	fmt.Println("\tEdits the todo with 33 and sets the due date to this coming Monday\n")
	yellow.Println("\ttodo e 33 due none")
	fmt.Println("\tEdits the todo with 33 and removes the due date\n")

	blueBold.Println("\nExpanding existing todos")
	yellow.Println("\ttodo ex 39 +final: read physics due mon, do literature report due fri")
	fmt.Println("\tRemoves the todo with id 39, and adds following two todos\n")

	blueBold.Println("\nDeleting")
	yellow.Println("\ttodo d 33")
	fmt.Println("\tDeletes a todo with id 33\n")
	fmt.Println("Todolist was lovingly crafted by Grant Ammons (https://twitter.com/gammons).")
	fmt.Println("For full documentation, please visit http://todolist.site")
}

func routeInput(command string, input string) {
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
	case "e", "edit":
		app.EditTodoDue(input)
	case "ex", "expand":
		app.ExpandTodo(input)
	case "init":
		app.InitializeRepo()
	case "web":
		if err := app.Load(); err != nil {
			os.Exit(1)
		} else {
			web := todolist.NewWebapp()
			fmt.Println("Now serving todolist web.\nHead to http://localhost:7890 to see your todo list!")
			open.Start("http://localhost:7890")
			web.Run()
		}
	}
}
