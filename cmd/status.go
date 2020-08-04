package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		deleteStatus        bool
		setStatusCmdDesc    = "Sets the status of a todo item"
		setStatusCmdExample = `
To set a status:
ultralist status 33 blocked
ultralist s 33 blocked

To remove a status:
ultralist status --remove 33`
		setStatusCmdLongDesc = `Sets the status of a todo item.  Status can be any string.`
	)

	var setStatusCmd = &cobra.Command{
		Use:     "status [id] <status>",
		Aliases: []string{"s"},
		Example: setStatusCmdExample,
		Long:    setStatusCmdLongDesc,
		Short:   setStatusCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			if deleteStatus {
				ultralist.NewApp().RemoveTodoStatus(strings.Join(args, " "))
			} else {
				ultralist.NewApp().SetTodoStatus(strings.Join(args, " "))
			}
		},
	}

	rootCmd.AddCommand(setStatusCmd)

	setStatusCmd.Flags().BoolVarP(&deleteStatus, "remove", "", false, "Removes status from a todo")
}
