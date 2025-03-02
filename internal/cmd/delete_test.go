package cmd

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/savabush/taskTracker/internal/services"
	"github.com/spf13/cobra"
)

func TestDeleteCmd_Args(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "No args",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "Empty arg",
			args:    []string{""},
			wantErr: true,
		},
		{
			name:    "Whitespace arg",
			args:    []string{"   "},
			wantErr: true,
		},
		{
			name:    "Valid task name",
			args:    []string{"Test task"},
			wantErr: false,
		},
		{
			name:    "Multiple args",
			args:    []string{"Task1", "Task2"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := DeleteCmd.Args(cmd, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteCmd_Run(t *testing.T) {
	var logBuf bytes.Buffer
	handler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	oldLogger := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(oldLogger)

	cleanup := createTempTaskFile(t)
	defer cleanup()

	tests := []struct {
		name           string
		setupTasks     []string
		taskToDelete   string
		expectedErr    bool
		expectedOutput string
	}{
		{
			name:           "Delete existing task",
			setupTasks:     []string{"Task1", "Task2", "Task3"},
			taskToDelete:   "Task2",
			expectedErr:    false,
			expectedOutput: "Deleted tasks",
		},
		{
			name:           "Delete non-existent task",
			setupTasks:     []string{"Task1", "Task3"},
			taskToDelete:   "NonExistentTask",
			expectedErr:    true,
			expectedOutput: "Deleted tasks", // The command actually doesn't log an error for this case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear log buffer
			logBuf.Reset()

			// Setup tasks
			service := services.NewTaskService()
			// Clear existing tasks
			for k := range service.Tasks {
				delete(service.Tasks, k)
			}
			// Add setup tasks
			for _, taskTitle := range tt.setupTasks {
				service.AddTask(taskTitle)
			}
			service.SaveTasks()

			// Get starting task count
			startingCount := len(service.Tasks)

			// Create a command
			cmd := &cobra.Command{}

			// Execute delete command
			DeleteCmd.Run(cmd, []string{tt.taskToDelete})

			// Reload service to see changes
			service = services.NewTaskService()

			// Check task count
			if !tt.expectedErr {
				if len(service.Tasks) != startingCount-1 {
					t.Errorf("Expected %d tasks after deletion, got %d", startingCount-1, len(service.Tasks))
				}

				// Verify the specific task was deleted
				for _, task := range service.Tasks {
					if task.Title == tt.taskToDelete {
						t.Errorf("Task '%s' was not deleted", tt.taskToDelete)
					}
				}
			} else {
				// For error cases, task count should remain the same
				if len(service.Tasks) != len(tt.setupTasks) {
					t.Errorf("Expected %d tasks to remain, got %d", len(tt.setupTasks), len(service.Tasks))
				}
			}

			// Check log output
			logOutput := logBuf.String()
			if !strings.Contains(logOutput, tt.expectedOutput) {
				t.Errorf("Expected log message '%s', but got: %s", tt.expectedOutput, logOutput)
			}
		})
	}
}
