package cmd

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/savabush/taskTracker/internal/services"
	"github.com/spf13/cobra"
)

func TestListCmd_Args(t *testing.T) {
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
			name:    "Valid filter - pending",
			args:    []string{"pending"},
			wantErr: false,
		},
		{
			name:    "Valid filter - inProgress",
			args:    []string{"inProgress"},
			wantErr: false,
		},
		{
			name:    "Valid filter - completed",
			args:    []string{"completed"},
			wantErr: false,
		},
		{
			name:    "Invalid filter",
			args:    []string{"invalid"},
			wantErr: true,
		},
		{
			name:    "Too many args",
			args:    []string{"pending", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := ListCmd.Args(cmd, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListCmd_Run(t *testing.T) {
	// Set up logging
	var logBuf bytes.Buffer
	handler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug})
	oldLogger := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(oldLogger)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Set up test file
	cleanup := createTempTaskFile(t)
	defer cleanup()

	// Prepare test data with different statuses
	service := services.NewTaskService()

	// Add pending task
	service.AddTask("Pending Task")

	// Add in-progress task
	service.AddTask("InProgress Task")
	service.InProgressTask("InProgress Task")

	// Add completed task
	service.AddTask("Completed Task")
	service.CompleteTask("Completed Task")

	service.SaveTasks()

	tests := []struct {
		name          string
		args          []string
		expectedTasks []string
	}{
		{
			name:          "List all tasks",
			args:          []string{},
			expectedTasks: []string{"Pending Task", "InProgress Task", "Completed Task"},
		},
		{
			name:          "List pending tasks",
			args:          []string{"pending"},
			expectedTasks: []string{"Pending Task"},
		},
		{
			name:          "List in-progress tasks",
			args:          []string{"inProgress"},
			expectedTasks: []string{"InProgress Task"},
		},
		{
			name:          "List completed tasks",
			args:          []string{"completed"},
			expectedTasks: []string{"Completed Task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset capture for each test
			if r != nil {
				r.Close()
			}
			r, w, _ = os.Pipe()
			os.Stdout = w

			// Clear log buffer
			logBuf.Reset()

			cmd := &cobra.Command{}
			ListCmd.Run(cmd, tt.args)

			// Close the pipe to capture output
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Check if all expected tasks are in the output
			for _, taskName := range tt.expectedTasks {
				if !strings.Contains(output, taskName) {
					t.Errorf("Expected output to contain '%s', but got: %s", taskName, output)
				}
			}

			// Check for debug logging
			if !strings.Contains(logBuf.String(), "Running list command") {
				t.Errorf("Expected debug log 'Running list command', but was not found")
			}

			// When filtering, check for appropriate filter log
			if len(tt.args) > 0 {
				if !strings.Contains(logBuf.String(), "Filtering tasks") {
					t.Errorf("Expected debug log about filtering, but was not found")
				}
			}
		})
	}
}
