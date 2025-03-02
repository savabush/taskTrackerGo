package cmd

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/savabush/taskTracker/internal/services"
	"github.com/spf13/cobra"
)

func createTempTaskFile(t *testing.T) func() {
	tmpFile, err := os.CreateTemp("", "tasks_test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	tmpFile.Write([]byte(`{"tasks":{}}`))
	tmpFile.Close()

	origFileName := services.GetTasksFileName()

	services.SetTasksFileName(tmpFile.Name())

	return func() {
		services.SetTasksFileName(origFileName)
		os.Remove(tmpFile.Name())
	}
}

func TestAddCmd_Args(t *testing.T) {
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
			name:    "Valid single arg",
			args:    []string{"Test task"},
			wantErr: false,
		},
		{
			name:    "Multiple args with one empty",
			args:    []string{"Task1", "   ", "Task3"},
			wantErr: true,
		},
		{
			name:    "Multiple valid args",
			args:    []string{"Task1", "Task2", "Task3"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := AddCmd.Args(cmd, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddCmd_Run(t *testing.T) {
	var logBuf bytes.Buffer
	handler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	oldLogger := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(oldLogger)

	cleanup := createTempTaskFile(t)
	defer cleanup()

	tests := []struct {
		name        string
		args        []string
		wantTaskCnt int
	}{
		{
			name:        "Add single task",
			args:        []string{"Test task"},
			wantTaskCnt: 1,
		},
		{
			name:        "Add multiple tasks",
			args:        []string{"Task1", "Task2", "Task3"},
			wantTaskCnt: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := services.NewTaskService()
			for k := range service.Tasks {
				delete(service.Tasks, k)
			}
			service.SaveTasks()

			cmd := &cobra.Command{}

			AddCmd.Run(cmd, tt.args)

			service = services.NewTaskService() // Reload from file
			count := 0
			for _, task := range service.Tasks {
				for _, arg := range tt.args {
					if task.Title == arg {
						count++
						break
					}
				}
			}

			if count != tt.wantTaskCnt {
				t.Errorf("Expected %d tasks to be added, but got %d", tt.wantTaskCnt, count)
			}

			logOutput := logBuf.String()
			if !strings.Contains(logOutput, "Added tasks") {
				t.Errorf("Expected log message 'Added tasks', but got: %s", logOutput)
			}
		})
	}
}

func TestAddCmd_Integration(t *testing.T) {
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	var logBuf bytes.Buffer
	handler := slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelInfo})
	oldLogger := slog.Default()
	slog.SetDefault(slog.New(handler))

	cleanup := createTempTaskFile(t)
	defer cleanup()
	defer func() {
		os.Stdout = oldStdout
		slog.SetDefault(oldLogger)
	}()

	AddCmd.SetArgs([]string{"Integration Test Task"})
	AddCmd.Execute()

	w.Close()

	service := services.NewTaskService()
	found := false
	for _, task := range service.Tasks {
		if task.Title == "Integration Test Task" {
			found = true
			if task.Status != services.TaskStatusPending {
				t.Errorf("Expected task status to be 'pending', but got '%s'", task.Status)
			}
			break
		}
	}

	if !found {
		t.Error("Task was not added correctly")
	}
}
