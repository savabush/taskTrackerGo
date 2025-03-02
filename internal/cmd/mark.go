package cmd

import (
	"errors"
	"log/slog"

	"github.com/savabush/taskTracker/internal/services"
	"github.com/spf13/cobra"
)

var MarkCompletedCmd = &cobra.Command{
	Use:   "mark-completed [task]",
	Short: "Mark a task as completed",
	Long:  `mark-completed is used to mark a task as completed.`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("filter must be one of: pending, in progress, completed")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		taskService := services.NewTaskService()
		_, err := taskService.GetTask(args[0])
		if err != nil {
			slog.Error("Task not found", "error", err)
			return
		}
		taskService.CompleteTask(args[0])
		taskService.SaveTasks()
		slog.Info("Marked task as completed", "task", args[0])
	},
}

var MarkInProgressCmd = &cobra.Command{
	Use:   "mark-in-progress [task]",
	Short: "Mark a task as in progress",
	Long:  `mark-in-progress is used to mark a task as in progress.`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("filter must be one of: pending, in progress, completed")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		taskService := services.NewTaskService()
		_, err := taskService.GetTask(args[0])
		if err != nil {
			slog.Error("Task not found", "error", err)
			return
		}
		taskService.InProgressTask(args[0])
		taskService.SaveTasks()
		slog.Info("Marked task as in progress", "task", args[0])
	},
}
