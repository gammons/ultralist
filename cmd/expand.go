package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	expandCmdDesc     = "Expands a project with todos"
	expandCmdExample  = "  ultralist expand 33 +importantProject: Get, Things, Done, ..."
	expandCmdLongDesc = expandCmdDesc + "."
)

var expandCmd = &cobra.Command{
	Use:     "expand [todo_id] +[project]: [todo_1], [todo_2], [todo_3], ...",
	Aliases: []string{"ex"},
	Example: expandCmdExample,
	Long:    expandCmdLongDesc,
	Short:   expandCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ExpandTodo("expand " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(expandCmd)
}
