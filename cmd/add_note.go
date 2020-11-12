package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		addNoteCmdDesc    = "Adds notes to todos.  Syntax is `an <todoID> <noteContent>`"
		addNoteCmdExample = `  ultralist an 1 this is a note for the first todo`
	)

	var addNoteCmd = &cobra.Command{
		Use:     "addnote",
		Aliases: []string{"an"},
		Example: addNoteCmdExample,
		Short:   addNoteCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			todoID, _ := strconv.Atoi(args[0])
			ultralist.NewApp().AddNote(todoID, strings.Join(args[1:], " "))
		},
	}
	rootCmd.AddCommand(addNoteCmd)
}
