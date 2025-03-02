package services

import (
	"os"
	"testing"
	"time"
)

func createTempTaskFile(t *testing.T) (string, func()) {
	tmpFile, err := os.CreateTemp("", "tasks_test*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	tmpFile.Write([]byte(`{"tasks":{}}`))
	tmpFile.Close()

	origFileName := GetTasksFileName()
	SetTasksFileName(tmpFile.Name())

	return tmpFile.Name(), func() {
		SetTasksFileName(origFileName)
		os.Remove(tmpFile.Name())
	}
}

func TestNewTaskService(t *testing.T) {
	// Set up a temporary file
	_, cleanup := createTempTaskFile(t)
	defer cleanup()

	// Create a new service
	service := NewTaskService()

	// Verify it's initialized correctly
	if service.Tasks == nil {
		t.Errorf("Expected Tasks map to be initialized, but it was nil")
	}
}

func TestAddTask(t *testing.T) {
	// Set up a temporary file
	_, cleanup := createTempTaskFile(t)
	defer cleanup()

	// Create a new service
	service := NewTaskService()

	// Test adding a task
	taskTitle := "Test Task"
	err := service.AddTask(taskTitle)
	if err != nil {
		t.Errorf("AddTask returned unexpected error: %v", err)
	}

	// Verify task was added
	if len(service.Tasks) != 1 {
		t.Errorf("Expected 1 task after adding, got %d", len(service.Tasks))
	}

	// Verify task properties
	task, exists := service.Tasks[taskTitle]
	if !exists {
		t.Errorf("Task '%s' was not found in the map", taskTitle)
	} else {
		if task.Title != taskTitle {
			t.Errorf("Task title doesn't match: got %s, want %s", task.Title, taskTitle)
		}
		if task.Status != TaskStatusPending {
			t.Errorf("New task should have pending status, got %s", task.Status)
		}
		if task.ID == "" {
			t.Errorf("Task ID should not be empty")
		}
		if task.CreatedAt.IsZero() {
			t.Errorf("CreatedAt should not be zero time")
		}
		if task.UpdatedAt.IsZero() {
			t.Errorf("UpdatedAt should not be zero time")
		}
	}
}

func TestGetTasks(t *testing.T) {
	// Set up a temporary file
	_, cleanup := createTempTaskFile(t)
	defer cleanup()

	// Create a new service with various tasks
	service := NewTaskService()

	// Add tasks with different statuses
	service.AddTask("Pending Task")

	service.AddTask("InProgress Task")
	service.InProgressTask("InProgress Task")

	service.AddTask("Completed Task")
	service.CompleteTask("Completed Task")

	// Test cases
	tests := []struct {
		name          string
		filter        TaskStatus
		expectedCount int
	}{
		{
			name:          "All tasks",
			filter:        "",
			expectedCount: 3,
		},
		{
			name:          "Pending tasks",
			filter:        TaskStatusPending,
			expectedCount: 1,
		},
		{
			name:          "In progress tasks",
			filter:        TaskStatusInProgress,
			expectedCount: 1,
		},
		{
			name:          "Completed tasks",
			filter:        TaskStatusCompleted,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filteredTasks := service.GetTasks(tt.filter)
			if len(filteredTasks) != tt.expectedCount {
				t.Errorf("Expected %d tasks with filter '%s', got %d",
					tt.expectedCount, tt.filter, len(filteredTasks))
			}

			// Check that tasks have correct status if filtering
			if tt.filter != "" {
				for _, task := range filteredTasks {
					if task.Status != tt.filter {
						t.Errorf("Task '%s' has incorrect status: got %s, want %s",
							task.Title, task.Status, tt.filter)
					}
				}
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	// Set up a temporary file
	_, cleanup := createTempTaskFile(t)
	defer cleanup()

	// Create a service and add a task
	service := NewTaskService()
	taskTitle := "Test Task"
	service.AddTask(taskTitle)

	// Test getting existing task
	task, err := service.GetTask(taskTitle)
	if err != nil {
		t.Errorf("Unexpected error getting existing task: %v", err)
	}
	if task.Title != taskTitle {
		t.Errorf("Task title doesn't match: got %s, want %s", task.Title, taskTitle)
	}

	// Test getting non-existent task
	_, err = service.GetTask("Non-existent Task")
	if err == nil {
		t.Errorf("Expected error getting non-existent task, but got nil")
	}
}

func TestDeleteTask(t *testing.T) {
	// Set up a temporary file
	_, cleanup := createTempTaskFile(t)
	defer cleanup()

	// Create a service and add a task
	service := NewTaskService()
	taskTitle := "Test Task"
	service.AddTask(taskTitle)

	// Test deleting existing task
	err := service.DeleteTask(taskTitle)
	if err != nil {
		t.Errorf("Unexpected error deleting task: %v", err)
	}
	if len(service.Tasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(service.Tasks))
	}

	// Test deleting non-existent task
	err = service.DeleteTask("Non-existent Task")
	if err == nil {
		t.Errorf("Expected error deleting non-existent task, but got nil")
	}
}

func TestSaveAndLoadTasks(t *testing.T) {
	// Set up a temporary file
	tmpFile, cleanup := createTempTaskFile(t)
	defer cleanup()

	// Create a service and add some tasks
	service := NewTaskService()
	service.AddTask("Task1")
	service.AddTask("Task2")
	service.CompleteTask("Task2")

	// Save tasks
	err := service.SaveTasks()
	if err != nil {
		t.Errorf("Unexpected error saving tasks: %v", err)
	}

	// Create a new service to load tasks
	service2 := NewTaskService()

	// Verify tasks were loaded correctly
	if len(service2.Tasks) != 2 {
		t.Errorf("Expected 2 tasks after loading, got %d", len(service2.Tasks))
	}

	// Verify task properties were preserved
	task, err := service2.GetTask("Task2")
	if err != nil {
		t.Errorf("Could not find task 'Task2' after reload")
	} else if task.Status != TaskStatusCompleted {
		t.Errorf("Task status not preserved: got %s, want %s", task.Status, TaskStatusCompleted)
	}

	// Test loading from empty file
	os.Remove(tmpFile)
	os.WriteFile(tmpFile, []byte{}, 0644)

	service3 := NewTaskService()
	if len(service3.Tasks) != 0 {
		t.Errorf("Expected 0 tasks when loading from empty file, got %d", len(service3.Tasks))
	}
}

func TestMarkTaskStatus(t *testing.T) {
	// Set up a temporary file
	_, cleanup := createTempTaskFile(t)
	defer cleanup()

	// Create a service and add a task
	service := NewTaskService()
	taskTitle := "Test Task"
	service.AddTask(taskTitle)

	// Record original updated time
	originalTime := service.Tasks[taskTitle].UpdatedAt

	// Wait a short time to ensure timestamps differ
	time.Sleep(10 * time.Millisecond)

	// Test marking as in progress
	err := service.InProgressTask(taskTitle)
	if err != nil {
		t.Errorf("InProgressTask returned unexpected error: %v", err)
	}

	task := service.Tasks[taskTitle]
	if task.Status != TaskStatusInProgress {
		t.Errorf("Expected task status to be in progress, got %s", task.Status)
	}
	if !task.UpdatedAt.After(originalTime) {
		t.Errorf("Expected UpdatedAt time to be updated")
	}

	// Test marking non-existent task
	err = service.InProgressTask("Non-existent Task")
	if err == nil {
		t.Errorf("Expected error marking non-existent task as in progress")
	}

	// Test marking as completed
	originalTime = service.Tasks[taskTitle].UpdatedAt
	time.Sleep(10 * time.Millisecond)

	err = service.CompleteTask(taskTitle)
	if err != nil {
		t.Errorf("CompleteTask returned unexpected error: %v", err)
	}

	task = service.Tasks[taskTitle]
	if task.Status != TaskStatusCompleted {
		t.Errorf("Expected task status to be completed, got %s", task.Status)
	}
	if !task.UpdatedAt.After(originalTime) {
		t.Errorf("Expected UpdatedAt time to be updated")
	}

	// Test marking non-existent task
	err = service.CompleteTask("Non-existent Task")
	if err == nil {
		t.Errorf("Expected error marking non-existent task as completed")
	}
}
