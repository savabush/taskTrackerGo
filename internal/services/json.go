package services

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	basePath = "data"
	fileName = "tasks.json"
)

func GetTasksFileName() string {
	return fileName
}

var tasksFileName = fileName

func SetTasksFileName(name string) {
	tasksFileName = name
}

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "inProgress"
	TaskStatusCompleted  TaskStatus = "completed"
)

var baseDataInData = []byte(`{"tasks":{}}`)

func createFileIfNotExists(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		slog.Debug("Creating new tasks file", "filename", filename)
		return os.WriteFile(filename, baseDataInData, 0644)
	}
	return nil
}

type Task struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type TaskService struct {
	Tasks map[string]Task `json:"tasks"`
	mu    sync.RWMutex
}

func NewTaskService() *TaskService {
	service := &TaskService{
		Tasks: make(map[string]Task),
	}
	service.LoadTasks()
	return service
}

func (s *TaskService) AddTask(title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	task := Task{
		ID:        uuid.New().String(),
		Title:     title,
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.Tasks[task.Title] = task
	return nil
}

func (s *TaskService) GetTasks(filter TaskStatus) map[string]Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	filteredTasks := make(map[string]Task)
	switch filter {
	case TaskStatusPending:
		for _, task := range s.Tasks {
			if task.Status == TaskStatusPending {
				filteredTasks[task.Title] = task
			}
		}
		return filteredTasks
	case TaskStatusInProgress:
		for _, task := range s.Tasks {
			if task.Status == TaskStatusInProgress {
				filteredTasks[task.Title] = task
			}
		}
		return filteredTasks
	case TaskStatusCompleted:
		for _, task := range s.Tasks {
			if task.Status == TaskStatusCompleted {
				filteredTasks[task.Title] = task
			}
		}
		return filteredTasks
	default:
		return s.Tasks
	}
}

func (s *TaskService) GetTask(title string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.Tasks[title]
	if !ok {
		return Task{}, errors.New("task not found")
	}
	return task, nil
}

func (s *TaskService) DeleteTask(title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.Tasks[title]
	if !ok {
		return errors.New("task not found")
	}
	delete(s.Tasks, title)
	return nil
}

func (s *TaskService) SaveTasks() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Debug("Saving tasks to file", "filename", tasksFileName, "count", len(s.Tasks))
	type TasksWrapper struct {
		Tasks map[string]Task `json:"tasks"`
	}

	wrapper := TasksWrapper{
		Tasks: s.Tasks,
	}

	jsonData, err := json.Marshal(wrapper)
	if err != nil {
		slog.Error("Failed to marshal tasks", "error", err)
		return errors.New("failed to marshal tasks")
	}
	slog.Debug("Writing tasks to file", "bytes", len(jsonData))
	return os.WriteFile(tasksFileName, jsonData, 0644)
}

func (s *TaskService) LoadTasks() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	slog.Debug("Loading tasks from file", "filename", tasksFileName)
	jsonData, err := os.ReadFile(tasksFileName)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Debug("Tasks file does not exist, creating new file")
			err = createFileIfNotExists(tasksFileName)
			if err != nil {
				slog.Error("Failed to create tasks file", "error", err)
				return err
			}
			return nil
		} else {
			slog.Error("Failed to read tasks file", "error", err)
			return errors.New("failed to read tasks")
		}
	}

	if len(jsonData) == 0 {
		slog.Debug("Tasks file is empty")
		return nil
	}

	type TasksWrapper struct {
		Tasks map[string]Task `json:"tasks"`
	}

	var wrapper TasksWrapper
	wrapper.Tasks = make(map[string]Task)

	slog.Debug("Unmarshaling tasks data", "bytes", len(jsonData))
	err = json.Unmarshal(jsonData, &wrapper)
	if err != nil {
		slog.Debug("Failed to unmarshal with wrapper format, trying old format", "error", err)
		err = json.Unmarshal(jsonData, &s.Tasks)
		if err != nil {
			slog.Error("Failed to unmarshal tasks data", "error", err)
			return errors.New("failed to unmarshal tasks")
		}
		slog.Debug("Successfully unmarshaled using old format", "tasks", len(s.Tasks))
	} else {
		s.Tasks = wrapper.Tasks
		slog.Debug("Successfully unmarshaled tasks", "count", len(s.Tasks))
	}

	return nil
}

func (s *TaskService) CompleteTask(title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.Tasks[title]
	if !ok {
		return errors.New("task not found")
	}
	task.Status = TaskStatusCompleted
	task.UpdatedAt = time.Now()
	s.Tasks[title] = task
	return nil
}

func (s *TaskService) InProgressTask(title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.Tasks[title]
	if !ok {
		return errors.New("task not found")
	}
	task.Status = TaskStatusInProgress
	task.UpdatedAt = time.Now()
	s.Tasks[title] = task
	return nil
}
