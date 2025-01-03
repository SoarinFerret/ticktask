package cmd

import (
	"errors"
	"strconv"
	"strings"

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
	Args: func(cmd *cobra.Command, args []string) error {
		// if int flag is set, no args are allowed
		id, _ := cmd.Flags().GetInt("id")
		if id != -1 && len(args) > 0 {
			return errors.New("cannot use both ID flag and args")
		}

		// all args must be integers
		for _, arg := range args {
			_, err := strconv.Atoi(arg)
			if err != nil {
				return errors.New("all args must be integers")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		ids := []int{}

		id, _ := cmd.Flags().GetInt("id")
		if id != -1 {
			ids = append(ids, id)
		}

		for _, arg := range args {
			id, _ := strconv.Atoi(arg)
			ids = append(ids, id)
		}

		// if ids is empty, lets do it interactively
		if len(ids) == 0 {
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

			// convert selectedOptions to ids by parsing the number before the first space
			for _, option := range selectedOptions {
				id, _ := strconv.Atoi(option[:strings.Index(option, " ")])
				ids = append(ids, id)
			}

		}

		for _, id := range ids {
			task, err := todotxt.CompleteTask(id)
			if err != nil {
				pExit("Error completing task:", err)
			}

			pterm.Info.Println("Task completed:", task.Todo)
		}

		// commit to git
		err := git.CommitTodo()
		if err != nil {
			pExit("Error committing task completion: ", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(closeCmd)

	closeCmd.Flags().IntP("id", "i", -1, "ID of task to complete")
}
