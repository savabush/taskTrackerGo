package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/savabush/taskTracker/internal/cmd"
	"github.com/savabush/taskTracker/internal/utils"
	"github.com/spf13/cobra"
)

func main() {
	// Define verbose flag
	var verbose bool

	// Create the root command
	var rootCmd = &cobra.Command{
		Use:   "taskTracker",
		Short: "A simple task tracker",
		Long:  `taskTracker is a simple task tracker that allows you to add, edit, and delete tasks.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Configure logger based on verbose flag
			logLevel := slog.LevelInfo
			if verbose {
				logLevel = slog.LevelDebug
			}

			// Create a pretty handler with appropriate level
			handler := utils.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						a.Value = slog.StringValue(time.Now().Format(time.RFC3339))
					}
					return a
				},
			})

			// Set as default logger
			slog.SetDefault(slog.New(handler))

			if verbose {
				slog.Debug("Starting taskTracker in debug mode")
			}
		},
	}

	// Add verbose flag to root command
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose (debug) logging")

	// Add commands
	rootCmd.AddCommand(cmd.AddCmd, cmd.ListCmd, cmd.MarkInProgressCmd, cmd.MarkCompletedCmd, cmd.DeleteCmd)

	// Execute root command
	rootCmd.Execute()
}
