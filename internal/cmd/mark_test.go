package cmd

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/savabush/taskTracker/internal/services"
	"github.com/spf13/cobra"
)

func TestMarkCompletedCmd_Args(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "No args",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "Valid task",
			args:    []string{"Task1"},
			wantErr: false,
		},
		{
			name:    "Too many args",
			args:    []string{"Task1", "Task2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := MarkCompletedCmd.Args(cmd, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkInProgressCmd_Args(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "No args",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "Valid task",
			args:    []string{"Task1"},
			wantErr: false,
		},
		{
			name:    "Too many args",
			args:    []string{"Task1", "Task2"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := MarkInProgressCmd.Args(cmd, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarkCompletedCmd_Run(t *testing.T) {
	// Set up logging
	var logBuf bytes.Buffer
	handler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	oldLogger := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(oldLogger)

	// Set up test file
	cleanup := createTempTaskFile(t)
	defer cleanup()

	// Test cases
	tests := []struct {
		name         string
		taskName     string
		taskExists   bool
		expectedLog  string
		expectedTask bool
	}{
		{
			name:         "Mark existing task as completed",
			taskName:     "Task1",
			taskExists:   true,
			expectedLog:  "Marked task as completed",
			expectedTask: true,
		},
		{
			name:         "Mark non-existent task",
			taskName:     "NonExistentTask",
			taskExists:   false,
			expectedLog:  "Task not found",
			expectedTask: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear log buffer
			logBuf.Reset()

			// Setup task
			service := services.NewTaskService()
			// Clear existing tasks
			for k := range service.Tasks {
				delete(service.Tasks, k)
			}

			// Add test task if it should exist
			if tt.taskExists {
				service.AddTask(tt.taskName)
			}

			service.SaveTasks()

			// Execute command
			cmd := &cobra.Command{}
			MarkCompletedCmd.Run(cmd, []string{tt.taskName})

			// Reload service to see changes
			service = services.NewTaskService()

			// Check task status
			if tt.expectedTask {
				task, err := service.GetTask(tt.taskName)
				if err != nil {
					t.Errorf("Task '%s' not found after marking completed", tt.taskName)
				} else if task.Status != services.TaskStatusCompleted {
					t.Errorf("Expected task status to be 'completed', got '%s'", task.Status)
				}
			}

			// Check log output
			logOutput := logBuf.String()
			if !strings.Contains(logOutput, tt.expectedLog) {
				t.Errorf("Expected log '%s', but got: %s", tt.expectedLog, logOutput)
			}
		})
	}
}

func TestMarkInProgressCmd_Run(t *testing.T) {
	// Set up logging
	var logBuf bytes.Buffer
	handler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	oldLogger := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(oldLogger)

	// Set up test file
	cleanup := createTempTaskFile(t)
	defer cleanup()

	// Test cases
	tests := []struct {
		name         string
		taskName     string
		taskExists   bool
		expectedLog  string
		expectedTask bool
	}{
		{
			name:         "Mark existing task as in progress",
			taskName:     "Task1",
			taskExists:   true,
			expectedLog:  "Marked task as in progress",
			expectedTask: true,
		},
		{
			name:         "Mark non-existent task",
			taskName:     "NonExistentTask",
			taskExists:   false,
			expectedLog:  "Task not found",
			expectedTask: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear log buffer
			logBuf.Reset()

			// Setup task
			service := services.NewTaskService()
			// Clear existing tasks
			for k := range service.Tasks {
				delete(service.Tasks, k)
			}

			// Add test task if it should exist
			if tt.taskExists {
				service.AddTask(tt.taskName)
			}

			service.SaveTasks()

			// Execute command
			cmd := &cobra.Command{}
			MarkInProgressCmd.Run(cmd, []string{tt.taskName})

			// Reload service to see changes
			service = services.NewTaskService()

			// Check task status
			if tt.expectedTask {
				task, err := service.GetTask(tt.taskName)
				if err != nil {
					t.Errorf("Task '%s' not found after marking in progress", tt.taskName)
				} else if task.Status != services.TaskStatusInProgress {
					t.Errorf("Expected task status to be 'inProgress', got '%s'", task.Status)
				}
			}

			// Check log output
			logOutput := logBuf.String()
			if !strings.Contains(logOutput, tt.expectedLog) {
				t.Errorf("Expected log '%s', but got: %s", tt.expectedLog, logOutput)
			}
		})
	}
}
