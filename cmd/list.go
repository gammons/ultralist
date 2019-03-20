package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	listCmdDesc    = "Listing todos including filtering and grouping"
	listCmdExample = `Filtering by date:

  ulralist list due (tod|today|tom|tomorrow|overdue|this week|next week|last week|mon|tue|wed|thu|fri|sat|sun|none)

  ulralist list due tod
  Lists all todos due today.

  ulralist list due tom
  Lists all todos due tomorrow.

  ulralist list due mon
  Lists all todos due monday.

Grouping:
You can group todos by context or project.

  ulralist list by c
  Lists all todos grouped by context.

  ulralist list by p
  Lists all todos grouped by project.

Grouping and filtering:
Of course, you can combine grouping and filtering to get a nice formatted list.

  ulralist list by c due today
  Lists all todos due today grouped by context.

  ulralist list +project by c due this week
  Lists all todos due today for +project, grouped by context.

  ulralist list @frank by p due tom
  Lists all todos due tomorrow concerining @frank for +project, grouped by project.`
	listCmdLongDesc = listCmdDesc + "."
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Example: listCmdExample,
	Long:    listCmdLongDesc,
	Short:   listCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ListTodos(strings.Join(args, " "))
	},
}

var (
	listArchivedCmdDesc     = "Lists all archived todos"
	listArchivedCmdLongDesc = listArchivedCmdDesc + "."
)

var listArchivedCmd = &cobra.Command{
	Use:     "archived",
	Aliases: []string{"a", "ar"},
	Long:    listArchivedCmdLongDesc,
	Short:   listArchivedCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ListTodos("archived")
	},
}

var (
	listCompletedCmdDesc    = "Lists all completed todos"
	listCompletedCmdExample = `  ulralist list completed (tod|today)
  Lists all todos that were completed today.

  ulralist list completed this week
  Lists all todos that were completed this week.`
	listCompletedCmdLongDesc = listArchivedCmdDesc + "."
)

var listCompletedCmd = &cobra.Command{
	Use:     "completed",
	Aliases: []string{"c"},
	Long:    listCompletedCmdLongDesc,
	Short:   listCompletedCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ListTodos("completed " + strings.Join(args, " "))
	},
}

var (
	listNotesCmdDesc    = "Lists all todos including notes"
	listNotesCmdExample = `  ultralist list notes
  List all todos including notes.

  ultralist list notes 33
  Lists todo 33 including its notes.`
	listNotesCmdLongDesc = listNotesCmdDesc + "."
)

var listNotesCmd = &cobra.Command{
	Use:     "notes [todo_id]",
	Aliases: []string{"n"},
	Example: listNotesCmdExample,
	Long:    listNotesCmdLongDesc,
	Short:   listNotesCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			i, err := strconv.Atoi(args[0])
			if err == nil {
				ultralist.NewApp().HandleNotes("n " + strconv.Itoa(i))
			} else {
				ultralist.NewApp().ListTodos("ln " + strings.Join(args, " "))
			}
		} else {
			ultralist.NewApp().ListTodos("ln " + strings.Join(args, " "))
		}
	},
}

var (
	listPrioritizeCmdDesc     = "Lists all prioritized todos"
	listPrioritizeCmdLongDesc = listPrioritizeCmdDesc + "."
)

var listPrioritizeCmd = &cobra.Command{
	Use:     "prioritized",
	Aliases: []string{"p"},
	Long:    listPrioritizeCmdLongDesc,
	Short:   listPrioritizeCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ListTodos("prioritized " + strings.Join(args, " "))
	},
}

var (
	listOverdueCmdDesc     = "Lists all todos where the due date is in the past"
	listOverdueCmdLongDesc = listOverdueCmdDesc + "."
)

var listOverdueCmd = &cobra.Command{
	Use:     "overdue",
	Aliases: []string{"o"},
	Long:    listOverdueCmdLongDesc,
	Short:   listOverdueCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ListTodos("overdue " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listArchivedCmd)
	listCmd.AddCommand(listCompletedCmd)
	listCmd.AddCommand(listNotesCmd)
	listCmd.AddCommand(listOverdueCmd)
	listCmd.AddCommand(listPrioritizeCmd)
}
