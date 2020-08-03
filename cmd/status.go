package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		setStatusCmdDesc    = "Sets the status of a todo item"
		setStatusCmdExample = `
ultralist status 33 blocked
ultralist s 33 blocked`
		setStatusCmdLongDesc = `Sets the status of a todo item.  Status can be any string.`
	)

	var setStatusCmd = &cobra.Command{
		Use:     "status [id] <status>",
		Aliases: []string{"s"},
		Example: setStatusCmdExample,
		Long:    setStatusCmdLongDesc,
		Short:   setStatusCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().SetTodoStatus(strings.Join(args, " "))
		},
	}

	rootCmd.AddCommand(setStatusCmd)
}
