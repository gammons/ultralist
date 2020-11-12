package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		long = `Delete a note from a todo.
  For more info, see https://ultralist.io/docs/cli/managing_tasks/#notes-management`
		example = `  To see your todos with notes:
    ultralist list --notes

  To delete note 0 from todo 3:
    ultralist dn 3 0`
	)

	var deleteNoteCmd = &cobra.Command{
		Use:     "deletenote <todoID> <noteID>",
		Aliases: []string{"dn"},
		Long:    long,
		Example: example,
		Short:   "Delete a note from a todo.",
		Run: func(cmd *cobra.Command, args []string) {
			todoID, _ := strconv.Atoi(args[0])
			noteID, _ := strconv.Atoi(args[1])
			ultralist.NewApp().DeleteNote(todoID, noteID)
		},
	}

	rootCmd.AddCommand(deleteNoteCmd)
}
