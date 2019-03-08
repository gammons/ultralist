package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var addCmd = &cobra.Command{
	Use:     "add [subject] [due]",
	Aliases: []string{"a"},
	Short:   "Add a todo",
	Long: `Add a todo
The 'a' command adds todos.
You can also optionally specify a due date.
Specify a due date by putting 'due <date>' at the end, where <date> is in (tod|today|tom|tomorrow|mon|tue|wed|thu|fri|sat|sun)

Examples for adding a todo:
ultralist a Meeting with @bob about +importantPrject due today
ultralist a +work +verify did @john fix the build\?`,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().AddTodo(strings.Join(args, " "))
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
