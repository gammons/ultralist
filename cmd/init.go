package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		initCmdDesc     = "Initializes a new todo list in the current directory"
		initCmdLongDesc = initCmdDesc + "."
	)

	var initCmd = &cobra.Command{
		Use:   "init",
		Long:  initCmdLongDesc,
		Short: initCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().InitializeRepo()
		},
	}

	rootCmd.AddCommand(initCmd)
}
