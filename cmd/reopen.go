package cmd

import (
	"strconv"

	"github.com/pterm/pterm"
	"github.com/soarinferret/ticktask/internal/git"
	"github.com/soarinferret/ticktask/internal/todotxt"
	"github.com/spf13/cobra"
)

var reopenCmd = &cobra.Command{
	Use:     "reopen",
	Aliases: []string{"undone", "uncomplete"},
	Short:   "Reopen a completed task",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			pExit("Error getting task ID:", err)
		}

		if id == -1 {
			// get id from args
			if len(args) != 1 {
				pExit("Error: reopen requires a single argument", nil)
			}
			id, err = strconv.Atoi(args[0])
			if err != nil {
				pExit("Invalid task ID:", err)
			}
		}

		task, err := todotxt.ReopenTask(id)
		if err != nil {
			pExit("Error reopening task:", err)
		}

		// commit to git
		err = git.CommitTodo()

		if err != nil {
			pExit("Error committing reopened task: ", err)
		}

		pterm.Info.Println("Task reopened:", task.Todo)
	},
}

func init() {
	rootCmd.AddCommand(reopenCmd)

	reopenCmd.Flags().IntP("id", "i", -1, "ID of task to complete")
}
