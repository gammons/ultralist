package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	revertCompletedTodo bool
	completeCmdDesc     = "Completes a todo"
	completeCmdExample  = `  ultralist complete 33
  Completes todo with id 33.

  ultralist complete 33 --revert
  Uncompletes todo with id 33.`
	completeCmdLongDesc = completeCmdDesc + "."
)

var completeCmd = &cobra.Command{
	Use:     "complete [id]",
	Aliases: []string{"c"},
	Example: completeCmdExample,
	Long:    completeCmdLongDesc,
	Short:   completeCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if revertCompletedTodo {
			ultralist.NewApp().UncompleteTodo(strings.Join(args, " "))
		} else {
			ultralist.NewApp().CompleteTodo(strings.Join(args, " "))
		}
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
	completeCmd.Flags().BoolVarP(&revertCompletedTodo, "revert", "", false, "Uncompletes a completed todo")
}
