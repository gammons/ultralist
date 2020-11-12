package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		deleteCmdDesc    = "Deletes todos"
		deleteCmdExample = `  To delete a todo with ID 33:
    ultralist d 33
    ultralist delete 33

  Note, this will also free up the id of 33.`
		deleteCmdLongDesc = `Delete a todo with a specified ID.

  See the full docs at https://ultralist.io/docs/cli/managing_tasks`
	)

	var deleteCmd = &cobra.Command{
		Use:     "delete [id]",
		Aliases: []string{"d", "rm"},
		Example: deleteCmdExample,
		Long:    deleteCmdLongDesc,
		Short:   deleteCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().DeleteTodo(strings.Join(args, " "))
		},
	}

	rootCmd.AddCommand(deleteCmd)
}
