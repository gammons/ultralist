package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	manageCmdDesc = "Manage your list with a simple terminal UI"
)

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: manageCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().OpenManager()
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)
}
