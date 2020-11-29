package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	var (
		webCmdDesc     = "Open your list on ultralist.io"
		webCmdLongDesc = "\nIf your list is synced with ultralist.io, use this command to open your list with your web browser."
	)

	var webCmd = &cobra.Command{
		Use:   "web",
		Long:  webCmdLongDesc,
		Short: webCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			// ultralist.NewApp().OpenWeb()
		},
	}

	rootCmd.AddCommand(webCmd)
}
