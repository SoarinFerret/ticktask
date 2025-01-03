package cmd

import (
	"strconv"

	"github.com/pterm/pterm"
	"github.com/soarinferret/ticktask/internal/git"
	"github.com/soarinferret/ticktask/internal/todotxt"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove tasks / todo from the file",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			pExit("Error getting task ID:", err)
		}

		if id == -1 {
			// get id from args
			if len(args) != 1 {
				pExit("Error: complete requires a single argument", nil)
			}
			id, err = strconv.Atoi(args[0])
			if err != nil {
				pExit("Invalid task ID:", err)
			}
		}

		// check task exists
		task, err := todotxt.GetTask(id)
		if err != nil {
			pExit("Task not found:", err)
		}

		// confirm
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			pterm.Info.Println("Task to remove:", task.String())
			result, _ := pterm.DefaultInteractiveConfirm.Show("Are you sure you want to remove this task?")
			if !result {
				return
			}
		}

		err = todotxt.RemoveTask(id)
		if err != nil {
			pExit("Error removing task:", err)
		}

		// commit changes
		// commit to git
		err = git.CommitTodo()
		if err != nil {
			pExit("Error committing task removal: ", err)
		}
		pterm.Info.Println("Task removed")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().IntP("id", "i", -1, "ID of task to remove")

	removeCmd.Flags().BoolP("force", "f", false, "Don't ask for confirmation")
}
