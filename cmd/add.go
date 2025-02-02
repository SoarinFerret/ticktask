/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"
	"time"

	"github.com/soarinferret/ticktask/internal/git"
	"github.com/soarinferret/ticktask/internal/todotxt"
	"github.com/spf13/cobra"

	"github.com/KEINOS/go-todotxt/todo"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add tasks / todo to the database",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		priority, _ := cmd.Flags().GetString("priority")

		timeSpent, _ := cmd.Flags().GetString("time")
		complete, _ := cmd.Flags().GetBool("complete")
		timeDuration, err := time.ParseDuration(timeSpent)
		if err != nil {
			pExit("Error parsing time: ", err)

		}

		// add all arguments as a single task
		task, err := todo.ParseTask(strings.Join(args, " "))
		if err != nil {
			pExit("Error parsing task: ", err)

		}

		// add time spent to task if not already present
		if _, exists := task.AdditionalTags["time"]; !exists {
			if task.AdditionalTags == nil {
				task.AdditionalTags = make(map[string]string)
			}
			task.AdditionalTags["time"] = timeDuration.String()
		}

		// add priority
		if priority != "" {
			p, err := formatPriority(priority)
			if err != nil {
				pExit("Error parsing priority: ", err)
			}
			task.Priority = p
		}

		// mark as complete if flag is set
		if complete {
			task.Complete()
		}

		task, err = todotxt.AddTask(*task)
		if err != nil {
			pExit("Error adding task: ", err)
		}

		// commit to git
		err = git.CommitTodo()
		if err != nil {
			pExit("Error committing task: ", err)
		}

		// print task
		printListTable([]todo.Task{*task})
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringP("time", "t", "0m", "The time spent on the task")
	addCmd.Flags().StringP("priority", "p", "", "The priority of the task (A-Z)")
	addCmd.Flags().BoolP("complete", "c", false, "Mark the task as complete")

}
