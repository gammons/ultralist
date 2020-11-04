package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	addCmdDesc    = "Adds todos and notes"
	addCmdExample = `  ultralist add Prepare meeting notes about +importantProject for the meeting with @bob due today
  ultralist add Meeting with @bob about +importantProject due:today
  ultralist add +work +verify did @john fix the build? due:tom`
	addCmdLongDesc = `Adds todos and notes.

You can optionally specify a due date.
This can be done by by putting 'due:<date>' at the end, where <date> is in (tod|today|tom|tomorrow|mon|tue|wed|thu|fri|sat|sun|thisweek|nextweek).

Dates can also be explicit, using 3 characters for the month.  They can be written in 2 different formats:

ultralist a buy flowers for mom due:may12
ultralist get halloween candy due:31oct`
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Example: addCmdExample,
	Long:    addCmdLongDesc,
	Short:   addCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().AddTodo(strings.Join(args, " "))
	},
}

var (
	addNoteCmdDesc    = "Adds a note to a todo"
	addNoteCmdExample = `  ultralist add note 33 Don't forget to reserve a meeting room
  Adds a new note to todo with id 33.`
	addNoteCmdLongDesc = "\nIf a todo has extra detail (like a URL), add the note to that todo with this command."
)

var addNoteCmd = &cobra.Command{
	Use:     "note [todo_id]",
	Aliases: []string{"n"},
	Example: addNoteCmdExample,
	Long:    addNoteCmdLongDesc,
	Short:   addNoteCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().HandleNotes(strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addNoteCmd)
}
