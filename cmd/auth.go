package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		webAuthCmdDesc     = "Authenticates you against ultralist.io"
		webAuthCmdLongDesc = `
This will authenticate your local ultralist with ultralist.io.

  Syncing with ultralist.io conveys many benefits:
    - Real-time sync with other ultralist binaries on other computers
    - Manage your lists via the web at app.ultralist.io
    - Use Ultralist on your mobile phone
    - Any many others.

  Ultralist.io is a paid service.  For more information, See https://ultralist.io/docs/cli/pro_integration`
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

	rootCmd.AddCommand(webAuthCmd)
	webAuthCmd.AddCommand(webAuthCheckCmd)
}
