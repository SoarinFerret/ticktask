package cmd

import (
	"errors"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/soarinferret/ticktask/internal/git"
	"github.com/soarinferret/ticktask/internal/todotxt"
	"github.com/spf13/cobra"
)

var priorityCmd = &cobra.Command{
	Use:     "priority",
	Aliases: []string{"p"},
	Short:   "Manage task priorities",
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

		// if 2 args, first must be an int & second must be a priority
		if len(args) == 2 {
			_, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New("first argument must be an integer")
			}

			_, err = formatPriority(args[1])
			if err != nil {
				return err
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// usage: tt priority 1 A
		// or
		// tt priority -i 1 -p A

		id, _ := cmd.Flags().GetInt("id")
		if id == -1 {
			id, _ = strconv.Atoi(args[0])
		}

		priority, _ := cmd.Flags().GetString("priority")
		if priority == "" {
			priority = args[1]
		}

		p, err := formatPriority(priority)
		if err != nil {
			pExit("Error parsing priority: ", err)
		}

		task, err := todotxt.SetPriority(id, p)
		if err != nil {
			pExit("Error setting priority: ", err)
		}

		pterm.Info.Println("Priority set:", task.String())

		// commit to git
		err = git.CommitTodo()
		if err != nil {
			pExit("Error committing task: ", err)
		}

	},
}

func formatPriority(priority string) (string, error) {
	priority = strings.ToUpper(priority)
	p := priority[0]
	if p >= 'A' && p <= 'Z' {
		return string(p), nil
	}
	return "", errors.New("priority must be a single uppercase letter")
}

func init() {
	rootCmd.AddCommand(priorityCmd)

	priorityCmd.Flags().IntP("id", "i", -1, "ID of task to set priority")
	priorityCmd.Flags().StringP("priority", "p", "", "Priority to set")
}
