package utils

import (
	"context"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	slog.Handler
	writer io.Writer
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	// Add colors based on log level
	levelText := level
	switch r.Level {
	case slog.LevelDebug:
		levelText = color.CyanString("%-7s", "DEBUG")
	case slog.LevelInfo:
		levelText = color.GreenString("%-7s", "INFO")
	case slog.LevelWarn:
		levelText = color.YellowString("%-7s", "WARN")
	case slog.LevelError:
		levelText = color.RedString("%-7s", "ERROR")
	}

	// Format timestamp
	timeStr := r.Time.Format("15:04:05.000")
	timeStr = color.MagentaString(timeStr)

	// Get the message
	msg := color.WhiteString(r.Message)

	// Write the log line
	_, err := io.WriteString(h.writer, timeStr+" "+levelText+" "+msg+"\n")
	return err
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(w, opts),
		writer:  w,
	}
}
