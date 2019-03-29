package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	webAuthCmdDesc     = "Authenticates you against ultralist.io"
	webAuthCmdLongDesc = `
Run "ultralist auth" to login or signup to ultralist.io to begin syncing lists.
"ultralist auth" will redirect you to Ultralist.io's login/signup page.
Ultralist stores a JSON web token in ~/.config/ultralist/creds.json.
`
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
	webAuthCheckCmdLongDesc = "\nCheck your login status to ultralist.io using \"ultralist auth check\"."
)

var webAuthCheckCmd = &cobra.Command{
	Use:   "check",
	Long:  webAuthCheckCmdLongDesc,
	Short: webAuthCheckCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().CheckAuth()
	},
}

func init() {
	rootCmd.AddCommand(webAuthCmd)
	webAuthCmd.AddCommand(webAuthCheckCmd)
}
