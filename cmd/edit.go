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
		editCmdDesc = "Edits todos"
		longDesc    = `Edits todos.

  You can edit all facets of a todo.

  Read the full docs at https://ultralist.io/docs/cli/managing_tasks/#editing-todos`
		editCmdExample = `  To edit a todo's subject:
    ultralist edit 33 Meeting with @bob about +project
    ultralist e 33 Change the subject once again

  To edit just the due date, keeping the subject:
    ultralist edit 33 due:mon

  To remove a due date:
    ultralist edit 33 due none

  To edit a status
    ultralist edit 33 status:next

	To remove a status:
    ultralist edit 33 status:none`
	)

	var editCmd = &cobra.Command{
		Use:     "edit [id]",
		Aliases: []string{"e"},
		Example: editCmdExample,
		Long:    longDesc,
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
