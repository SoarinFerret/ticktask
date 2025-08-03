/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"regexp"
	"strconv"
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

		date, _ := cmd.Flags().GetString("date")
		y, _ := cmd.Flags().GetBool("yesterday")
		if y {
			date = "-1d"
		}

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

		// parse date if provided
		if date != "0d" {
			t, err := customDateParser(date)
			if err != nil {
				pExit("Error parsing date: ", err)
			}
			task.CreatedDate = t
			if task.Completed {
				task.CompletedDate = t
			}
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
	addCmd.Flags().StringP("date", "d", "0d", "Date the task should have been added (e.g. 2024-01-01 or -1d for yesterday)")
	addCmd.Flags().BoolP("yesterday", "y", false, "Shortcut for --date=-1d")
}

func customDateParser(date string) (time.Time, error) {
	// if date is in the format of +1d or -1d, parse it as a duration. Format accepts days, weeks, and months
	if strings.HasPrefix(date, "+") || strings.HasPrefix(date, "-") {
		// custom duration parser using regexp
		days, err := parseCustomDuration(date)
		if err != nil {
			return time.Time{}, err
		}

		return time.Now().AddDate(0, 0, days), nil
	}
	// otherwise, parse it as a date
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func parseCustomDuration(input string) (days int, err error) {
	re := regexp.MustCompile(`^([+-])(\d+)([dwmy])$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	sign, valueStr, unit := matches[1], matches[2], matches[3]
	value, _ := strconv.Atoi(valueStr)

	// Convert unit to days
	switch unit {
	case "d":
		// 1 day
	case "w":
		value *= 7
	case "m":
		value *= 30
	case "y":
		value *= 365
	default:
		return 0, fmt.Errorf("unknown unit: %s", unit)
	}

	if sign == "-" {
		value = -value
	}

	return value, nil
}
