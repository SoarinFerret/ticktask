package cmd

import (
	"errors"
	"strconv"
	"time"

	"github.com/pterm/pterm"
	"github.com/soarinferret/ticktask/internal/git"
	"github.com/soarinferret/ticktask/internal/todotxt"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"l"},
	Short:   "Log time spent on a task",
	Args: func(cmd *cobra.Command, args []string) error {
		// args must be 0 or 2
		if len(args) != 0 && len(args) != 2 {
			return errors.New("requires 0 or 2 arguments")
		}

		// if no args, flags must be set
		if len(args) == 0 {
			if !cmd.Flags().Changed("id") || !cmd.Flags().Changed("time") {
				return errors.New("requires ID and time flags if no arguments are present")
			}
		}

		// if 2 args, first must be an int & second must be a duration
		if len(args) == 2 {
			_, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New("first argument must be an integer")
			}

			_, err = time.ParseDuration(args[1])
			if err != nil {
				return errors.New("second argument must be a duration")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// usage: tt log -i 1 -t 1h
		// or
		// tt log 1 1h

		// get id from flags or args
		id, _ := cmd.Flags().GetInt("id")
		if id == -1 {
			id, _ = strconv.Atoi(args[0])
		}

		// get time from flags or args
		timeSpent, _ := cmd.Flags().GetString("time")
		if timeSpent == "" {
			timeSpent = args[1]
		}

		timeDuration, err := time.ParseDuration(timeSpent)
		if err != nil {
			pExit("Error parsing time: ", err)
		}

		override, _ := cmd.Flags().GetBool("override")

		task, err := todotxt.AddTimeToTask(id, timeDuration, override)
		if err != nil {
			pExit("Error adding time to task:", err)
		}

		// commit to git
		err = git.CommitTodo()
		if err != nil {
			pExit("Error committing time log: ", err)
		}

		// print new task time
		pterm.Info.Println("Time updated for", strconv.Itoa(task.ID)+":", task.AdditionalTags["time"])
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().IntP("id", "i", -1, "ID of task to complete")
	logCmd.Flags().StringP("time", "t", "", "Time spent on task")
	logCmd.Flags().BoolP("override", "o", false, "Override existing time log entry")
}
