package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var deleteNoteCmd = &cobra.Command{
		Use:     "deletenote",
		Aliases: []string{"dn"},
		Example: "Deletes notes from todos.  Syntax is `dn <todoID> <noteID>`",
		Short:   "Delete a note from a todo",
		Run: func(cmd *cobra.Command, args []string) {
			todoID, _ := strconv.Atoi(args[0])
			noteID, _ := strconv.Atoi(args[1])
			ultralist.NewApp().DeleteNote(todoID, noteID)
		},
	}

	rootCmd.AddCommand(deleteNoteCmd)
}
