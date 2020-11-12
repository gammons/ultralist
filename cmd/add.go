package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		addCmdDesc    = "Adds todos"
		addCmdExample = `  ultralist add Prepare meeting notes about +importantProject for the meeting with @bob due today
		ultralist add Meeting with @bob about +importantProject due:today
		ultralist add +work +verify did @john fix the build? due:tom`
		addCmdLongDesc = `Adds todos.

	You can optionally specify a due date.
	This can be done by by putting 'due:<date>' at the end, where <date> is in (tod|today|tom|tomorrow|mon|tue|wed|thu|fri|sat|sun|thisweek|nextweek).

	Dates can also be explicit, using 3 characters for the month.  They can be written in 2 different formats:

	ultralist a buy flowers for mom due:may12
	ultralist get halloween candy due:31oct`
	)

	var addCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Example: addCmdExample,
		Long:    addCmdLongDesc,
		Short:   addCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().AddTodo(strings.Join(args, " "))
		},
	}

	rootCmd.AddCommand(addCmd)
}
