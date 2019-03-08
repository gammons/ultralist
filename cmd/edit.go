package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Args:    cobra.ArbitraryArgs,
	Aliases: []string{"e"},
	Short:   "edit a todo",
	Long: `Edit a todo
The 'edit' or 'e' command will edit a todo.
You can edit a todo's due date or subject.

Example - change due date to this coming Monday, and preserve existing subject:
ultralist e 33 due mon

Change the subject, preseve due date:
ultralist e 33 meet with @bob for lunch

Change both the subject and the due date:
ultralist e 33 meet with @bob for lunch due mon`,

	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().EditTodo(strings.Join(args, " "))
	},
}

func init() {
	RootCmd.AddCommand(editCmd)
}
