package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	editCmdDesc    = "Edits todos and notes"
	editCmdExample = `  ultralist edit 33 Meeting with @bob about +importantProject and +anotherImportantProject due today
  Edits todo 33 entirely.

  ultralist edit 33 due mon
  Edits todo 33 and sets the due date to next Monday.

  ultralist edit 33 due none
  Edits todo 33 and removes the due date.`
	editCmdLongDesc = editCmdDesc + "."
)

var editCmd = &cobra.Command{
	Use:     "edit [id]",
	Aliases: []string{"e"},
	Example: editCmdExample,
	Long:    editCmdLongDesc,
	Short:   editCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().EditTodo("edit " + strings.Join(args, " "))
	},
}

var (
	editNoteCmdDesc    = "Edits a note from a todo"
	editNoteCmdExample = `  ultralist edit note 33 3 Don't forget to reserve a meeting room with phone
  Edits the 3rd note of todo with id 33.`
	editNoteCmdLongDesc = addNoteCmdDesc + "."
)

var editNoteCmd = &cobra.Command{
	Use:     "note [todo_id] [id]",
	Aliases: []string{"n"},
	Example: editNoteCmdExample,
	Long:    editNoteCmdLongDesc,
	Short:   editNoteCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().HandleNotes("en " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.AddCommand(editNoteCmd)
}
