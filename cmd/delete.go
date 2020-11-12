package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		deleteCmdDesc    = "Deletes todos"
		deleteCmdExample = `  ultralist delete 33
		Deletes todo with id 33.`
		deleteCmdLongDesc = deleteCmdDesc + "."
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
