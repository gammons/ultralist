package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		startCmdDesc    = "Start work on a task"
		startCmdExample = `  ultralist start 33
		ultralist s 33`
		startCmdLongDesc = `Starts work on a task.  Sets the StartedDate which can then be filtered upon.`

		stopCmdDesc     = "Stop work on a started task"
		stopCmdExample  = `  ultralist stop 33`
		stopCmdLongDesc = `Stops work on a task.  Un-sets the StartedDate.`
	)

	var startCmd = &cobra.Command{
		Use:     "start [id]",
		Aliases: []string{"s"},
		Example: startCmdExample,
		Long:    startCmdLongDesc,
		Short:   startCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().StartTodo(strings.Join(args, " "))
		},
	}

	var stopCmd = &cobra.Command{
		Use:     "stop [id]",
		Example: stopCmdExample,
		Long:    stopCmdLongDesc,
		Short:   stopCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().StopTodo(strings.Join(args, " "))
		},
	}

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}
