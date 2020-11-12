package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		editCmdDesc    = "Edits todos"
		editCmdExample = `  ultralist edit 33 Meeting with @bob about +importantProject and +anotherImportantProject due today
		Edits todo 33 entirely.

		ultralist edit 33 due mon
		Edits todo 33 and sets the due date to next Monday.

		ultralist edit 33 due none
		Edits todo 33 and removes the due date.`
		editCmdLongDesc = editCmdDesc + "."
	)

	var editCmd = &cobra.Command{
		Use:     "edit [id]",
		Aliases: []string{"e"},
		Example: editCmdExample,
		Long:    editCmdLongDesc,
		Short:   editCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			todoID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Could not parse todo ID: '%s'\n", args[0])
				return
			}
			ultralist.NewApp().EditTodo(todoID, strings.Join(args[1:], " "))
		},
	}

	rootCmd.AddCommand(editCmd)
}
