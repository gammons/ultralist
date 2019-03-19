package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	webCmdDesc     = "Authenticate and synchronize against ultralist.io"
	webCmdLongDesc = webCmdDesc
)

var webCmd = &cobra.Command{
	Use:   "web",
	Long:  webCmdLongDesc,
	Short: webCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().OpenWeb()
	},
}

var (
	webAuthCmdDesc     = "Authenticates you against ultralist.io"
	webAuthCmdLongDesc = webAuthCmdDesc
)

var webAuthCmd = &cobra.Command{
	Use:   "auth",
	Long:  webAuthCmdLongDesc,
	Short: webAuthCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().AuthWorkflow()
	},
}

var (
	webAuthCheckCmdDesc     = "Checks your authentication status against ultralist.io"
	webAuthCheckCmdLongDesc = webAuthCheckCmdDesc
)

var webAuthCheckCmd = &cobra.Command{
	Use:   "check",
	Long:  webAuthCheckCmdLongDesc,
	Short: webAuthCheckCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().CheckAuth()
	},
}

var (
	webSyncCmdDesc     = "Syncs todos with ultralist.io"
	webSyncCmdLongDesc = webSyncCmdDesc
)

var webSyncCmd = &cobra.Command{
	Use:   "sync",
	Long:  webSyncCmdLongDesc,
	Short: webSyncCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().Sync(strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.AddCommand(webAuthCmd)
	webCmd.AddCommand(webSyncCmd)
	webAuthCmd.AddCommand(webAuthCheckCmd)
}
