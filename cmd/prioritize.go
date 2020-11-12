package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		prioritizeCmdExample = `
		ultralist prioritize 33
		ultralist p 33
			Prioritizes todo with id 33.

		ultralist unprioritize 33
		ultralist up 33
			Unprioritizes todo with id 33.`
		prioritizeCmdLongDesc = `Prioritize and unprioritize todos.

	Todos with a priority flag will be highlighted on top of the list.`
	)

	var prioritizeCmd = &cobra.Command{
		Use:     "prioritize [id]",
		Aliases: []string{"p"},
		Example: prioritizeCmdExample,
		Long:    prioritizeCmdLongDesc,
		Short:   "Prioritize a todo",
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().PrioritizeTodo(strings.Join(args, " "))
		},
	}

	var unprioritizeCmd = &cobra.Command{
		Use:     "unprioritize [id]",
		Aliases: []string{"up"},
		Example: prioritizeCmdExample,
		Long:    prioritizeCmdLongDesc,
		Short:   "Un-prioritize a todo",
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().UnprioritizeTodo(strings.Join(args, " "))
		},
	}

	rootCmd.AddCommand(prioritizeCmd)
	rootCmd.AddCommand(unprioritizeCmd)
}
