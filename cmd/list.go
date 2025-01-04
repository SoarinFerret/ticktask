package cmd

import (
	"strconv"
	"time"

	"github.com/KEINOS/go-todotxt/todo"
	"github.com/spf13/cobra"

	"github.com/pterm/pterm"
	"github.com/soarinferret/ticktask/internal/todotxt"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all tasks",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

		projects, contexts := globalArgsHandler(args)

		list, err := todotxt.GetTasks()

		for _, project := range projects {
			list = list.Filter(todo.FilterByProject(project))
		}
		for _, context := range contexts {
			list = list.Filter(todo.FilterByContext(context))
		}

		if err != nil {
			pExit("Error loading task list: ", err)
		}

		list.Sort(todo.SortPriorityAsc)

		printListTable(list)
	},
}

var listTodoCmd = &cobra.Command{
	Use:     "todo",
	Aliases: []string{"t"},
	Short:   "List all incomplete tasks",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		projects, contexts := globalArgsHandler(args)

		list, err := todotxt.GetTasks()

		list = list.Filter(todo.FilterNotCompleted)

		for _, project := range projects {
			list = list.Filter(todo.FilterByProject(project))
		}
		for _, context := range contexts {
			list = list.Filter(todo.FilterByContext(context))
		}

		if err != nil {
			pExit("Error loading task list: ", err)
		}

		list.Sort(todo.SortDueDateAsc)
		list.Sort(todo.SortPriorityAsc)

		printListTable(list)

		return
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(listTodoCmd)

	// listCmd Filters
	//listCmd.Flags().BoolP()
}

func printListTable(list todo.TaskList) {

	listData := [][]string{}

	// check if tasks include due date
	dueDate := false
	for _, task := range list {
		if task.DueDate != (time.Time{}) {
			dueDate = true
			break
		}
	}

	// check if tasks include completion date
	completionDate := false
	for _, task := range list {
		if task.CompletedDate != (time.Time{}) {
			completionDate = true
			break
		}
	}

	// Add header
	if dueDate && completionDate {
		listData = append(listData, []string{"ID", "Due", "Pri", "Task", "TimeSpent", "Completed"})
	} else if dueDate {
		listData = append(listData, []string{"ID", "Due", "Pri", "Task", "TimeSpent"})
	} else if completionDate {
		listData = append(listData, []string{"ID", "Pri", "Task", "TimeSpent", "Completed"})
	} else {
		listData = append(listData, []string{"ID", "Pri", "Task", "TimeSpent"})
	}

	for _, task := range list {
		dueDateStr := ""
		if task.DueDate == (time.Time{}) {
			dueDateStr = "-"
		} else {
			dueDateStr = task.DueDate.Format("2006-01-02")
		}

		priorityStr := "-"
		if task.Priority != "" {
			priorityStr = task.Priority
		}

		completedStr := "-"
		if task.CompletedDate != (time.Time{}) {
			completedStr = task.CompletedDate.Format("2006-01-02")
		}

		listItem := []string{}

		if dueDate && completionDate {
			listItem = []string{
				strconv.Itoa(task.ID),
				dueDateStr,
				priorityStr,
				task.Todo,
				task.AdditionalTags["time"],
				completedStr,
			}
		} else if dueDate {
			listItem = []string{
				strconv.Itoa(task.ID),
				dueDateStr,
				priorityStr,
				task.Todo,
				task.AdditionalTags["time"],
			}
		} else if completionDate {
			listItem = []string{
				strconv.Itoa(task.ID),
				priorityStr,
				task.Todo,
				task.AdditionalTags["time"],
				completedStr,
			}
		} else {
			listItem = []string{
				strconv.Itoa(task.ID),
				priorityStr,
				task.Todo,
				task.AdditionalTags["time"],
			}
		}

		if dueDate && task.IsOverdue() {
			// highlight overdue tasks red
			listItem[1] = pterm.LightRed(listItem[1])
		} else if dueDate && task.IsDueToday() {
			// highlight due today tasks yellow
			listItem[1] = pterm.LightYellow(listItem[1])
		}

		listData = append(listData, listItem)
	}

	pterm.DefaultTable.WithHasHeader().WithData(listData).Render()
}
