package cmd

import (
	//"strconv"

	//"strings"

	"fmt"
	"sort"
	"time"

	"github.com/KEINOS/go-todotxt/todo"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/soarinferret/ticktask/internal/report"
	"github.com/soarinferret/ticktask/internal/todotxt"
)

var reportCmd = &cobra.Command{
	Use:     "report",
	Aliases: []string{"r"},
	Short:   "Generate task summary reports",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

		t, _ := cmd.Flags().GetString("time")

		projects, contexts := globalArgsHandler(args)

		list, err := todotxt.GetTasks()
		if err != nil {
			pExit("Error loading task list: ", err)
		}
		list.Filter(todo.FilterCompleted)

		// Filter tasks by time range
		startTime, err := report.CalculateOldest(t)
		if err != nil {
			pExit("Error calculating start time for report: ", err)
		}
		endTime := time.Now()
		list = list.Filter(todotxt.FilterByDateRange(startTime, endTime))

		for _, project := range projects {
			list = list.Filter(todo.FilterByProject(project))
		}
		for _, context := range contexts {
			list = list.Filter(todo.FilterByContext(context))
		}

		list.Sort(todo.SortCompletedDateDesc)

		summary, err := report.Generate(list, t)
		if err != nil {
			pExit("Error generating report: ", err)
		}

		printReportTable(summary)

	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringP("time", "t", "7d", "Time range for the report (e.g., 7d, 4w, 3m)")

	// TODO: reportCmd Filters
	//reportCmd.Flags().BoolP("by-context", "c", false, "List tasks by context instead of date")
	//reportCmd.Flags().BoolP("by-project", "p", false, "List tasks by project instead of date")
}

func printReportTable(summary map[string]report.Summary) {

	listData := [][]string{}

	// Add header
	listData = append(listData, []string{"Date", "Time Spent", "Tasks Completed"})

	// Sort the summary by date
	var dates []string
	for date := range summary {
		dates = append(dates, date)
	}
	// Sort dates in descending order
	sort.Slice(dates, func(i, j int) bool {
		t1, _ := time.Parse("2006-01-02", dates[i])
		t2, _ := time.Parse("2006-01-02", dates[j])
		return t2.After(t1)
	})
	for _, date := range dates {
		data := summary[date]
		// Add each date's summary data to the list
		listData = append(listData, []string{
			date,
			data.TimeSpent.String(),
			fmt.Sprintf("%d", data.Count),
		})
	}

	pterm.DefaultTable.WithHasHeader().WithHeaderRowSeparator("-").WithData(listData).Render()
}
