package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		archiveCompletedTodo bool
		long                 = `Completes or un-completes a todo.

For more info, see https://ultralist.io/docs/cli/managing_tasks`
		completeCmdExample = `
  Complete a todo with id 33:
    ultralist complete 33
    ultralist c 33

  Complete a todo with id 33 and archive it:
    ultralist uncomplete 33 --archive

  Uncompletes todo with id 33.
    ultralist uncomplete 33
    ultralist uc 33`
	)

	var completeCmd = &cobra.Command{
		Use:     "complete [id]",
		Aliases: []string{"c"},
		Example: completeCmdExample,
		Short:   "Completes a todo.",
		Long:    long,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().CompleteTodo(strings.Join(args, " "), archiveCompletedTodo)
		},
	}

	var uncompleteCmd = &cobra.Command{
		Use:     "uncomplete [id]",
		Aliases: []string{"uc"},
		Example: completeCmdExample,
		Short:   "Un-completes a todo.",
		Long:    long,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().UncompleteTodo(strings.Join(args, " "))
		},
	}

	rootCmd.AddCommand(completeCmd)
	completeCmd.Flags().BoolVarP(&archiveCompletedTodo, "archive", "", false, "Archives a completed todo automatically")
	rootCmd.AddCommand(uncompleteCmd)
}
