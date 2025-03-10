package cmd

import (
	"errors"
	"log/slog"
	"strings"
	"sync"

	"github.com/savabush/taskTracker/internal/services"
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [task]",
	Short: "Delete a task",
	Long:  `delete is used to delete a task from the task list`,
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

		wg := sync.WaitGroup{}
		wg.Add(len(args))

		for _, arg := range args {
			go func(arg string) {
				defer wg.Done()
				taskService.DeleteTask(arg)
			}(arg)
		}
		wg.Wait()
		taskService.SaveTasks()
		slog.Info("Deleted tasks", "count", len(args))
	},
}
