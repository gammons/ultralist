package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	deleteCmdDesc    = "Deletes todos and notes"
	deleteCmdExample = `  ultralist delete 33
  Deletes todo with id 33.`
	deleteCmdLongDesc = deleteCmdDesc
)

var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Aliases: []string{"d", "rm"},
	Example: deleteCmdExample,
	Long:    deleteCmdLongDesc,
	Short:   deleteCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().DeleteTodo(strings.Join(args, " "))
	},
}

var (
	deleteNoteCmdDesc    = "Deletes a note from a todo"
	deleteNoteCmdExample = `  ultralist delete note 33 3
  Deletes the 3rd note of the todo with id 33.`
	deleteNoteCmdLongDesc = addNoteCmdDesc
)

var deleteNoteCmd = &cobra.Command{
	Use:     "note [todo_id] [id]",
	Aliases: []string{"n"},
	Example: deleteNoteCmdExample,
	Long:    deleteNoteCmdLongDesc,
	Short:   deleteNoteCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().HandleNotes("dn " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteNoteCmd)
}
