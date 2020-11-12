package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		setStatusCmdDesc    = "Sets the status of a todo item"
		setStatusCmdExample = `  To add a "blocked" status to a todo:
    ultralist status 33 blocked
    ultralist s 33 blocked

  You can remove a status by setting a status to "none".  Example:
    ultralist s 33 none`

		setStatusCmdLongDesc = `Sets the status of a todo item.
  A status should be a single lower-case word, e.g. "now", "blocked", or "waiting".

  For more info, see https://ultralist.io/docs/cli/managing_tasks/#handling-todo-statuses`
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
