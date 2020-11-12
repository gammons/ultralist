package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		syncCmdDesc    = "Sync a list with ultralist.io"
		syncCmdExample = `  To synchronize your list:
    ultralist sync

  To set up your list to sync with ultralist.io:
    ultralist sync --setup

  To stop syncing your list with ultralist.io:
    ultralist sync --unsync

  Perform a sync without showing output to the screen:
    ultralist sync --quiet
	`

		syncCmdQuiet    bool
		setupCmd        bool
		unsyncCmd       bool
		syncCmdLongDesc = `Sync a list with ultralist.io.
  If you're using Ultralist Pro, this will manually (and bi-directionally) sync your list with the remote list on ultralist.io.

  Note that you won't normally have to run this command.  Syncing occurs automatically when you manipulate your list locally.

  For more info on syncing, see https://ultralist.io/docs/cli/pro_integration`
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
