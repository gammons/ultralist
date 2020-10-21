package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	archiveCompletedTodo bool
	completeCmdExample   = `
  ultralist complete 33
  ultralist c 33
    Completes todo with id 33.

  ultralist uncomplete 33 --archive
    Completes todo with id 33 and archives it.

  ultralist uncomplete 33
	ultralist uc 33
    Uncompletes todo with id 33.`
)

var completeCmd = &cobra.Command{
	Use:     "complete [id]",
	Aliases: []string{"c"},
	Example: completeCmdExample,
	Short:   "Completes a todo.",
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().CompleteTodo(strings.Join(args, " "), archiveCompletedTodo)
	},
}

var uncompleteCmd = &cobra.Command{
	Use:     "uncomplete [id]",
	Aliases: []string{"uc"},
	Example: completeCmdExample,
	Short:   "Un-completes a todo.",
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().UncompleteTodo(strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
	completeCmd.Flags().BoolVarP(&archiveCompletedTodo, "archive", "", false, "Archives a completed todo automatically")
	rootCmd.AddCommand(uncompleteCmd)
}
