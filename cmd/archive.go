package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

func init() {
	var (
		archiveCmdExample = `
  To arvhive a todo with id 33:
    ultralist archive 33
    ultralist ar 33

  To unarchive todo with id 33:
    ultralist unarchive 33
    ultralist uar 33

  To archive all completed todos:
    ultralist archive completed
    ultralist ar c

  Garbage collection will delete all archived todos, reclaming ids:
    ultralist archive gc
    ultralist ar gc

  ultralist archive gc
  ultralist ar gc
	  Run garbage collection. Delete all archived todos and reclaim ids

  See the full docs here:
  https://ultralist.io/docs/cli/managing_tasks`
	)

	var archiveCmd = &cobra.Command{
		Use:     `archive [id]`,
		Aliases: []string{"ar"},
		Example: archiveCmdExample,
		Short:   "Archives a todo.",
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().ArchiveTodo(strings.Join(args, " "))
		},
	}

	var unarchiveCmd = &cobra.Command{
		Use:     "unarchive [id]",
		Aliases: []string{"uar"},
		Example: archiveCmdExample,
		Short:   "Un-archives a todo.",
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().UnarchiveTodo(strings.Join(args, " "))
		},
	}

	var archiveCompletedCmd = &cobra.Command{
		Use:     "c",
		Example: "  ultralist archive completed\n  ultralist ar c",
		Short:   "Achives all completed todos.",
		Long: `Achives all completed todos.
For more info, see https://ultralist.io/docs/cli/managing_tasks/#archivingunarchiving-todos`,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().ArchiveCompleted()
		},
	}

	var archiveGarbageCollectCmd = &cobra.Command{
		Use:     "gc",
		Aliases: []string{"rm"},
		Short:   "Deletes all archived todos.",
		Long: `Delete all archived todos, reclaiming ids.
For more info, see https://ultralist.io/docs/cli/managing_tasks/#archivingunarchiving-todos`,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewApp().GarbageCollect()
		},
	}

	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(unarchiveCmd)
	archiveCmd.AddCommand(archiveCompletedCmd)
	archiveCmd.AddCommand(archiveGarbageCollectCmd)
}
