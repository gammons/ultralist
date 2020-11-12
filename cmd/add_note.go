package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		addNoteCmdDesc     = "Adds notes to todos."
		addNoteLongCmdDesc = addNoteCmdDesc + "\n For more info, see https://ultralist.io/docs/cli/managing_tasks/#notes-management"
		addNoteCmdExample  = "  ultralist an 1 this is a note for the first todo"
	)

	var addNoteCmd = &cobra.Command{
		Use:     "addnote <todoID> <noteContent>",
		Aliases: []string{"an"},
		Example: addNoteCmdExample,
		Short:   addNoteCmdDesc,
		Long:    addNoteLongCmdDesc,
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			todoID, _ := strconv.Atoi(args[0])
			ultralist.NewApp().AddNote(todoID, strings.Join(args[1:], " "))
		},
	}
	rootCmd.AddCommand(addNoteCmd)
}
