package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	unicodeSupport bool
	colorSupport   bool
	listNotes      bool
	listCmdDesc    = "Listing todos including filtering and grouping"
	listCmdExample = `
Filtering by date:
------------------

  ultralist list due:(tod|today|tom|tomorrow|agenda|overdue|thisweek|nextweek|lastweek|mon|tue|wed|thu|fri|sat|sun|none)

  List all todos due today:
    ultralist list due:tod

  Lists all todos due tomorrow:
    ultralist list due:tom

  Lists all todos due monday:
    ultralist list due:mon

  Lists all overdue todos:
    ultralist list due:overdue

  Lists all todos whose due date is today or earlier:
    ultralist list due:agenda

Filtering by priority, completed, etc:
--------------------------------------

  You can filter todos on their priority or completed status:
    ultralist list is:priority
    ultralist list not:priority

    ultralist list is:completed
    ultralist list not:completed

  There are additional filters for showing completed todos:
    ultralist list completed:today
    ultralist list completed:thisweek

  By default, ultralist will not show archived todos. To show archived todos:
    ultralist list is:archived

Grouping:
---------
  You can group todos by context or project.

  Lists all todos grouped by context:
    ultralist list group:c

  Lists all todos grouped by project:
    ultralist list group:p

Combining filters:
-----------------------

  Of course, you can combine grouping and filtering to get a nice formatted list.

  Lists all todos due today grouped by context:
    ultralist list group:c due:today

  Lists all todos due today for +project, grouped by context:
    ultralist list +project group:c due:thisweek

  Lists all prioritized todos that are not completed and are overdue.  Include a todo's notes when listing:
    ultralist list --notes is:priority not:completed due:overdue

  Lists all todos due tomorrow concerining @frank for +project, grouped by project:
    ultralist list @frank group:p due:tom

For complete documentation, see https://ultralist.io/docs/cli/showing_tasks
`
	listCmdLongDesc = listCmdDesc + "."
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "ls"},
	Example: listCmdExample,
	Long:    listCmdLongDesc,
	Short:   listCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewAppWithPrintOptions(unicodeSupport, colorSupport).ListTodos(strings.Join(args, " "), listNotes)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&unicodeSupport, "unicode", "", true, "Allows unicode support in Ultralist output")
	listCmd.Flags().BoolVarP(&colorSupport, "color", "", true, "Allows color in Ultralist output")
	listCmd.Flags().BoolVarP(&listNotes, "notes", "", false, "Show a todo's notes when listing. ")
}
