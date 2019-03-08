package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var RootCmdBLUB string = fmt.Sprintf("Ultralist v%s, simple task management for tech folks.\n", ultralist.VERSION)

var RootCmd = &cobra.Command{
	Use:   "ultralist",
	Short: RootCmdBLUB,
	Long:  RootCmdBLUB,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.Usage()
	},
}

func Execute() {
	RootCmd.PersistentFlags().Bool("color", true, "use colors in output")
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
