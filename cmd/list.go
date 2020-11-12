package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		unicodeSupport bool
		colorSupport   bool
		listNotes      bool
		showStatus     bool
		listCmdDesc    = "List todos."
		listCmdExample = `
  Filtering by date:
  ------------------

  ultralist list due:<date>
  ultralist list duebefore:<date>
  ultralist list dueafter:<date>

  where <date> is one of:
  (tod|today|tom|tomorrow|thisweek|nextweek|lastweek|mon|tue|wed|thu|fri|sat|sun|none|<specific date>)

  List all todos due today:
    ultralist list due:tod

  Lists all todos due tomorrow:
    ultralist list due:tom

  Lists all todos due monday:
    ultralist list due:mon

  Lists all todos with no due date:
    ultralist list due:none

  Lists all overdue todos:
    ultralist list duebefore:today

  Lists all todos in due in the future:
    ultralist list dueafter:today

  When using a specific date, it needs to be in the format of jun23 or 23jun:
    ultralist list due:jun23

  Filtering by status:
  --------------------

  List all todos with a status of "started"
    ultralist list status:started

  List all todos without a status of "started"
    ultralist list status:-started

  List all todos without a status of "started" or "finished"
    ultralist list status:-started,-finished

  Filtering by projects or contexts:
  ----------------------------------

  Project and context filtering are very similar:
    ultralist list project:<project>
    ultralist list context:<context>

  List all todos with a project of "mobile"
    ultralist list project:mobile

  List all todos with a project of "mobile" and "devops"
    ultralist list project:mobile,devops

  List all todos with a project of "mobile" but not "devops"
    ultralist list project:mobile,-devops

  List all todos without a project of "devops"
    ultralist list project:-devops

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

  List all todos grouped by context:
    ultralist list group:c

  List all todos grouped by project:
    ultralist list group:p

  List all todos grouped by status:
 	  ultralist list group:s

  Combining filters:
  ------------------

  Of course, you can combine grouping and filtering to get a nice formatted list.

  Lists all todos due today grouped by context:
    ultralist list group:c due:today

  Lists all todos due today for +mobile, grouped by context:
    ultralist list project:mobile group:c due:thisweek

  Lists all prioritized todos that are not completed and are overdue.  Include a todo's notes when listing:
    ultralist list --notes is:priority not:completed duebefore:tod

  Lists all todos due tomorrow concerning @frank for +project, grouped by project:
    ultralist list context:frank group:p due:tom

  Indicator flags
  ---------------

  If you pass --status=true as a flag, you'll see an extra column when listing todos.

  * = Todo is prioritized
  N = Todo has notes attached
  A = Todo is archived
`
		listCmdLongDesc = `List todos, optionally providing a filter.

When listing todos, you can apply powerful filters, and perform grouping.

See the full docs at https://ultralist.io/docs/cli/showing_tasks`
	)

	var listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Example: listCmdExample,
		Long:    listCmdLongDesc,
		Short:   listCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewAppWithPrintOptions(unicodeSupport, colorSupport).ListTodos(strings.Join(args, " "), listNotes, showStatus)
		},
	}

	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&unicodeSupport, "unicode", "", true, "Allows unicode support in Ultralist output")
	listCmd.Flags().BoolVarP(&colorSupport, "color", "", true, "Allows color in Ultralist output")
	listCmd.Flags().BoolVarP(&listNotes, "notes", "", false, "Show a todo's notes when listing. ")
	listCmd.Flags().BoolVarP(&showStatus, "status", "", false, "Show a todo's status")
}
