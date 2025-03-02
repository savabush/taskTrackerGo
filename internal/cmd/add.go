package cmd

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/savabush/taskTracker/internal/services"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [# strings to add]",
	Short: "Add a new task",
	Long:  `add is used to add a new task to the task list`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a task description")
		}
		for _, arg := range args {
			if len(strings.TrimSpace(arg)) == 0 {
				return errors.New("task description cannot be empty")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		taskService := services.NewTaskService()
		for _, arg := range args {
			taskService.AddTask(arg)
		}
		taskService.SaveTasks()
		slog.Info("Added tasks", "count", len(args))
	},
}
