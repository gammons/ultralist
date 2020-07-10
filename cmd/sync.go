package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	syncCmdDesc     = "Sync a list with ultralist.io"
	syncCmdExample  = "ultralist sync"
	syncCmdQuiet    bool
	syncCmdLongDesc = `

The sync command has a few uses:
* If you are logged into ultralist.io (via "ultralist auth"), you can sync an existing list using "ultralist sync".
* For a synced list, you can run "ultralist sync" explicitly to receive any changes that may have occurred elsewhere.

Local changes to a synced list automatically get pushed to ultralist.io.

See https://ultralist.io/docs/cli/pro_integration for more info.
`
)

var syncCmd = &cobra.Command{
	Use:     "sync",
	Example: syncCmdExample,
	Long:    syncCmdLongDesc,
	Short:   syncCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().Sync(syncCmdQuiet)
	},
}

func init() {
	syncCmd.Flags().BoolVarP(&syncCmdQuiet, "quiet", "q", false, "Run without output")
	rootCmd.AddCommand(syncCmd)
}
