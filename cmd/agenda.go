package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
)

var (
	agendaCmdDesc     = "Lists all todos where the due date is today or in the past"
	agendaCmdLongDesc = agendaCmdDesc + "."
)

var agendaCmd = &cobra.Command{
	Use:   "agenda",
	Long:  agendaCmdLongDesc,
	Short: agendaCmdDesc,
	Run: func(cmd *cobra.Command, args []string) {
		ultralist.NewApp().ListTodos("agenda " + strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(agendaCmd)
}
