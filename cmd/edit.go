package cmd

import (
	"github.com/pterm/pterm"
	"github.com/soarinferret/ticktask/internal/config"
	"github.com/soarinferret/ticktask/internal/edit"
	"github.com/soarinferret/ticktask/internal/git"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit the todo.txt file",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		edit.OpenFile(config.GetTodoTxtPath())

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
