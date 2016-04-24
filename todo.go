package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/gammons/todolist/todolist"
)

func main() {
	store := todolist.NewFileStore()
	store.Load()

	doOutput(store.Data)
}

func doOutput(todos []todolist.Todo) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for _, item := range todos {
		str := strconv.Itoa(item.Id) + "\t" + completedText(item.Completed) + "\t" + item.Subject
		fmt.Fprintln(w, str)
		//fmt.Fprintln(w,
	}
	w.Flush()
}

func completedText(completed bool) string {
	if completed {
		return "[x]"
	} else {
		return "[ ]"
	}
}
