package cmd

import (
	"strconv"

	"github.com/KEINOS/go-todotxt/todo"
	"github.com/pterm/pterm"
	"github.com/soarinferret/ticktask/internal/git"
	"github.com/soarinferret/ticktask/internal/todotxt"
	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:     "close",
	Aliases: []string{"done", "complete"},
	Short:   "Mark a task as complete",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			pExit("Error getting task ID:", err)
		}

		if id == -1 && len(args) == 1 {
			id, err = strconv.Atoi(args[0])
			if err != nil {
				pExit("Invalid task ID:", err)
			}
		}

		if id != -1 {
			task, err := todotxt.CompleteTask(id)
			if err != nil {
				pExit("Error completing task:", err)
			}

			// commit to git
			err = git.CommitTodo()

			if err != nil {
				pExit("Error committing task completion: ", err)
			}

			pterm.Info.Println("Task completed:", task.Todo)
			return
		}

		// else, lets do it interactively?
		tasks, err := todotxt.GetTasks()
		if err != nil {
			pExit("Error loading task list: ", err)
		}

		// filter out completed tasks
		tasks = tasks.Filter(todo.FilterNotCompleted)

		var options []string
		for _, task := range tasks {
			options = append(options, strconv.Itoa(task.ID)+" | "+task.Todo)
		}

		selectedOptions, _ := pterm.DefaultInteractiveMultiselect.WithOptions(options).Show()

		pterm.Info.Printfln("Selected options: %s", pterm.Green(selectedOptions))

	},
}

func init() {
	rootCmd.AddCommand(closeCmd)

	closeCmd.Flags().IntP("id", "i", -1, "ID of task to complete")
}
