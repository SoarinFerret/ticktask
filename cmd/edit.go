package cmd

import (
	"strconv"
	"github.com/pterm/pterm"
	"github.com/soarinferret/ticktask/internal/config"
	"github.com/soarinferret/ticktask/internal/edit"
	"github.com/soarinferret/ticktask/internal/git"
	"github.com/soarinferret/ticktask/internal/todotxt"
	"github.com/spf13/cobra"
	"github.com/KEINOS/go-todotxt/todo"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit the todo.txt file",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		// if no arguments, open the todo.txt file
		if len(args) == 0 {
			edit.OpenFile(config.GetTodoTxtPath())
		} else {
			id, _ := strconv.Atoi(args[0])

			task, err := todotxt.GetTask(id)
			if err != nil {
				pExit("Error getting task:", err)
			}

			taskString, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue(task.String()).Show()

			newTask, err := todo.ParseTask(taskString)
			if err != nil {
				pExit("Error parsing task:", err)
			}

			newTask.ID = task.ID
			task, err = todotxt.EditTask(*newTask)
			if err != nil {
				pExit("Error editing task:", err)
			}
		}

		noCommit, _ := cmd.Flags().GetBool("no-commit")
		if !noCommit {
			err := git.CommitTodo()
			if err != nil {
				pExit("Error committing manual edit:", err)
			}
			pterm.Info.Println("Committed Changes")
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().BoolP("no-commit", "n", false, "Don't commit changes to git")
}
