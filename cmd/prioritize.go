package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		example = `  To prioritize a todo with id 33:
    ultralist prioritize 33
    ultralist p 33

  To un-prioritize a todo with an id 33:
    ultralist unprioritize 33
    ultralist up 33`

		long = `Prioritize and un-prioritize todos.

  Todos with the priority flag will be highlighted, and will be at the top of your list.

  For more info, see https://ultralist.io/docs/cli/managing_tasks/#prioritizingunprioritizing-todos`
	)

	var prioritizeCmd = &cobra.Command{
		Use:     "prioritize [id]",
		Aliases: []string{"p"},
		Example: example,
		Long:    long,
		Short:   "Prioritize a todo.",
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().PrioritizeTodo(strings.Join(args, " "))
		},
	}

	var unprioritizeCmd = &cobra.Command{
		Use:     "unprioritize [id]",
		Aliases: []string{"up"},
		Example: example,
		Long:    long,
		Short:   "Un-prioritize a todo.",
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().UnprioritizeTodo(strings.Join(args, " "))
		},
	}

	rootCmd.AddCommand(prioritizeCmd)
	rootCmd.AddCommand(unprioritizeCmd)
}
