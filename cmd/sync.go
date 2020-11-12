package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		syncCmdDesc    = "Sync a list with ultralist.io"
		syncCmdExample = `ultralist sync
	ultralist sync --setup
	ultralist sync --unsync
	ultralist sync --quiet
	`

		syncCmdQuiet    bool
		setupCmd        bool
		unsyncCmd       bool
		syncCmdLongDesc = `Sync a list with Ultralist Pro.

	ultralist sync
		The Ultralist CLI stores tasks locally.  When running sync, it will push
		up any local changes to ultralist.io, as well as pulling down any remote changes from ultralist.io.

	ultralist sync --setup
		Set up the local list to sync with ultralist.io.  Or, pull a list from Ultralist.io to local.

	ultralist sync --unsync
		Stop syncing a local list with ultralist.io.

	ultralist sync --quiet
		Perform a sync without showing output to the screen.

	See https://ultralist.io/docs/cli/pro_integration for more info.
	`
	)

	var syncCmd = &cobra.Command{
		Use:     "sync",
		Example: syncCmdExample,
		Long:    syncCmdLongDesc,
		Short:   syncCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			if setupCmd {
				ultralist.NewApp().SetupSync()
				return
			}
			if unsyncCmd {
				ultralist.NewApp().Unsync()
				return
			}

			ultralist.NewApp().Sync(syncCmdQuiet)
		},
	}

	syncCmd.Flags().BoolVarP(&syncCmdQuiet, "quiet", "q", false, "Run without output")
	syncCmd.Flags().BoolVarP(&setupCmd, "setup", "", false, "Set up a list to sync with ultralist.io, or pull a remote list to local")
	syncCmd.Flags().BoolVarP(&unsyncCmd, "unsync", "", false, "Stop syncing a list with ultralist.io")
	rootCmd.AddCommand(syncCmd)
}
