package utils

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestPrettyHandler(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create a handler with the buffer
	handler := NewPrettyHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Create a logger with our handler
	logger := slog.New(handler)

	// Test different log levels
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// Get the output
	output := buf.String()

	// Verify each level appears in the output
	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Output doesn't contain DEBUG level")
	}
	if !strings.Contains(output, "INFO") {
		t.Errorf("Output doesn't contain INFO level")
	}
	if !strings.Contains(output, "WARN") {
		t.Errorf("Output doesn't contain WARN level")
	}
	if !strings.Contains(output, "ERROR") {
		t.Errorf("Output doesn't contain ERROR level")
	}

	// Verify messages appear
	if !strings.Contains(output, "Debug message") {
		t.Errorf("Output doesn't contain debug message")
	}
	if !strings.Contains(output, "Info message") {
		t.Errorf("Output doesn't contain info message")
	}
	if !strings.Contains(output, "Warning message") {
		t.Errorf("Output doesn't contain warning message")
	}
	if !strings.Contains(output, "Error message") {
		t.Errorf("Output doesn't contain error message")
	}
}

func TestPrettyHandlerLevelFiltering(t *testing.T) {
	// Create a buffer to capture the output
	var buf bytes.Buffer

	// Create a handler with info level threshold
	handler := NewPrettyHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Create a logger with our handler
	logger := slog.New(handler)

	// Log at different levels
	logger.Debug("Debug message") // Should be filtered out
	logger.Info("Info message")   // Should appear
	logger.Error("Error message") // Should appear

	// Get the output
	output := buf.String()

	// Debug message should not appear
	if strings.Contains(output, "Debug message") {
		t.Errorf("Output contains debug message when it should be filtered out")
	}

	// Info and error messages should appear
	if !strings.Contains(output, "Info message") {
		t.Errorf("Output doesn't contain info message")
	}
	if !strings.Contains(output, "Error message") {
		t.Errorf("Output doesn't contain error message")
	}
}

func TestPrettyHandlerFormat(t *testing.T) {
	var buf bytes.Buffer

	// Setup a record to test the formatting directly
	record := slog.Record{
		Time:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Message: "Test message",
		Level:   slog.LevelInfo,
	}

	// Create handler and test directly
	handler := NewPrettyHandler(&buf, nil)
	err := handler.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Handle returned unexpected error: %v", err)
	}

	output := buf.String()

	// Check for basic components in output
	if !strings.Contains(output, "Test message") {
		t.Errorf("Output doesn't contain the message")
	}

	if !strings.Contains(output, "INFO") {
		t.Errorf("Output doesn't contain the level")
	}

	// Verify output format (time, level, message)
	parts := strings.Fields(strings.TrimSpace(output))
	if len(parts) < 3 {
		t.Errorf("Output doesn't have expected format. Got: %s", output)
	}
}
