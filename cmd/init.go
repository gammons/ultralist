package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		initCmdDesc     = "Initializes a new todo list in the current directory"
		initCmdLongDesc = `Initializes a new todo list in the current directory.

This will create a .todos.json in the directory you're in.  You can then start adding todos to it.

For more info, see https://ultralist.io/docs/cli/managing_lists`
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
