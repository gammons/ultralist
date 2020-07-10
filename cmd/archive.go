package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	revertArchivedTodo bool
	archiveCmdDesc     = "Archives todos"
	archiveCmdExample  = `  ultralist archive 33
  Archives todo with id 33.

  ultralist archive 33 --revert
  Unarchives todo with id 33.`
	archiveCmdLongDesc = archiveCmdDesc + "."
)

var archiveCmd = &cobra.Command{
	Use:     "archive [id]",
	Aliases: []string{"ar"},
	Example: archiveCmdExample,
	Long:    archiveCmdLongDesc,
	Short:   archiveCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		if revertArchivedTodo {
			ultralist.NewApp().UnarchiveTodo(strings.Join(args, " "))
		} else {
			ultralist.NewApp().ArchiveTodo(strings.Join(args, " "))
		}
	},
}

var (
	archiveCompletedCmdDesc     = "Archives all completed todos"
	archiveCompletedCmdLongDesc = archiveCompletedCmdDesc + "."
)

var archiveCompletedCmd = &cobra.Command{
	Use:     "completed",
	Aliases: []string{"c"},
	Example: "ultralist archive completed",
	Long:    archiveCompletedCmdLongDesc,
	Short:   archiveCompletedCmdDesc,
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
	Long:    archiveGarbageCollectCmdLongDesc,
	Short:   archiveGarbageCollectCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().GarbageCollect()
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().BoolVarP(&revertArchivedTodo, "revert", "", false, "Unarchives an archived todo")
	archiveCmd.AddCommand(archiveCompletedCmd)
	archiveCmd.AddCommand(archiveGarbageCollectCmd)
}
