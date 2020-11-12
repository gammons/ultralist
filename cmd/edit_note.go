package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		cmdDesc     = "Edits notes on a todo."
		longCmdDesc = "Edits notes on a todo.\n For more info, see https://ultralist.io/docs/cli/managing_tasks/#notes-management"
		example     = `  To see your todos with notes:
    ultralist list --notes

  To edit note 0 from todo 3:
    ultralist en 3 0 this is the new note`
	)
	var editNoteCmd = &cobra.Command{
		Use:     "editnote",
		Aliases: []string{"en"},
		Example: example,
		Long:    longCmdDesc,
		Short:   cmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			todoID, _ := strconv.Atoi(args[0])
			noteID, _ := strconv.Atoi(args[1])
			ultralist.NewApp().EditNote(todoID, noteID, strings.Join(args[2:], " "))
		},
	}

	rootCmd.AddCommand(editNoteCmd)
}
