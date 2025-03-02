package cmd

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/savabush/taskTracker/internal/services"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list [filter]",
	Short: "List all tasks",
	Long:  `list is used to list all tasks. If a filter is provided, it will filter the tasks by the status. The filter must be one of: pending, in progress, completed`,

	Args: func(cmd *cobra.Command, args []string) error {
		slog.Debug("Validating list command arguments", "args", args)
		if len(args) > 1 {
			slog.Debug("Too many arguments provided")
			return errors.New("filter must be one of: pending, in progress, completed")
		}

		if len(args) == 1 {
			slog.Debug("Checking filter value", "filter", args[0])
			if args[0] != string(services.TaskStatusPending) && args[0] != string(services.TaskStatusInProgress) && args[0] != string(services.TaskStatusCompleted) {
				slog.Debug("Invalid filter value")
				return errors.New("filter must be one of: pending, in progress, completed")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("Running list command")
		taskService := services.NewTaskService()

		var filter services.TaskStatus
		if len(args) == 1 {
			filter = services.TaskStatus(args[0])
			slog.Debug("Filtering tasks", "filter", filter)
		} else {
			slog.Debug("No filter provided, showing all tasks")
		}

		tasks := taskService.GetTasks(filter)
		slog.Debug("Retrieved tasks from service", "count", len(tasks))

		for _, task := range tasks {
			fmt.Printf("%s %s - %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"), task.Title, task.Status)
		}
	},
}
