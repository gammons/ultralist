package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	archiveCmdDesc    = "Archives and un-archives todos"
	archiveCmdExample = `
  ultralist archive 33
  ultralist ar 33
    Archives todo with id 33.

  ultralist unarchive 33
  ultralist uar 33
    Unarchives todo with id 33.

  ultralist archive completed
  ultralist ar c
	  archives all completed todos

  ultralist archive gc
  ultralist ar gc
	  Run garbage collection. Delete all archived todos and reclaim ids`
)

var archiveCmd = &cobra.Command{
	Use:     "archive [id]",
	Aliases: []string{"ar"},
	Example: archiveCmdExample,
	Short:   "Archives a todo",
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ArchiveTodo(strings.Join(args, " "))
	},
}

var unarchiveCmd = &cobra.Command{
	Use:     "unarchive [id]",
	Aliases: []string{"aar"},
	Example: archiveCmdExample,
	Short:   "Un-archives a todo",
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().UnarchiveTodo(strings.Join(args, " "))
	},
}

var archiveCompletedCmd = &cobra.Command{
	Use:     "completed",
	Aliases: []string{"c"},
	Example: "ultralist archive completed",
	Short:   "Achives all completed todos",
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ArchiveCompleted()
	},
}

var (
	archiveGarbageCollectCmdDesc     = "Deletes all archived todos"
	archiveGarbageCollectCmdLongDesc = "\nDelete all archived todos, and reclaim ids"
)

var archiveGarbageCollectCmd = &cobra.Command{
	Use:     "garbage-collect",
	Aliases: []string{"gc", "rm"},
	Short:   "Deletes all archived todos",
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().GarbageCollect()
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(unarchiveCmd)
	archiveCmd.AddCommand(archiveCompletedCmd)
	archiveCmd.AddCommand(archiveGarbageCollectCmd)
}
