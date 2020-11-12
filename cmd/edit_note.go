package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var editNoteCmd = &cobra.Command{
		Use:     "editnote",
		Aliases: []string{"en"},
		Example: "edits notes from todos.  Syntax is `en <todoID> <noteID> <noteContent>`",
		Short:   "edit a note on a todo.",
		Run: func(cmd *cobra.Command, args []string) {
			todoID, _ := strconv.Atoi(args[0])
			noteID, _ := strconv.Atoi(args[1])
			ultralist.NewApp().EditNote(todoID, noteID, strings.Join(args[2:], " "))
		},
	}

	rootCmd.AddCommand(editNoteCmd)
}
