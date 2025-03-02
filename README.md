# Task Tracker

A simple CLI application to manage your tasks with ease, written in Go.

## Features

- **Task Management**: Add, list, delete, and update task status
- **Task Statuses**: Track tasks as pending, in-progress, or completed
- **Filtering**: List tasks by status
- **Colored Output**: Easy-to-read colorized terminal output
- **Logging**: Configurable logging with pretty formatting
- **Persistent Storage**: Tasks are saved between sessions

## Installation

### Prerequisites

- Go 1.18 or higher

### From Binary Release

You can download pre-built binaries for your platform from the [GitHub Releases](https://github.com/savabush/taskTracker/releases) page.

1. Download the archive for your platform
2. Extract the executable
3. Move the executable to a location in your PATH (optional)

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/savabush/taskTracker.git
   cd taskTracker
   ```

2. Build the application:
   ```
   go build -o task-tracker
   ```

3. Move the binary to your PATH (optional):
   ```
   mv task-tracker /usr/local/bin/
   ```

## Usage

### Adding Tasks

Add a new task to your list:

```
task-tracker add "Complete the project report"
```

### Listing Tasks

List all tasks:

```
task-tracker list
```

Filter tasks by status:

```
task-tracker list pending
task-tracker list inProgress
task-tracker list completed
```

### Updating Task Status

Mark a task as in-progress:

```
task-tracker mark-in-progress "Complete the project report"
```

Mark a task as completed:

```
task-tracker mark-completed "Complete the project report"
```

### Deleting Tasks

Delete a task:

```
task-tracker delete "Complete the project report"
```

### Verbose Logging

Enable detailed logs with the `--verbose` flag:

```
task-tracker --verbose list
```

## Code Structure

```
.
├── cmd/
│   └── task-tracker/        # Main entry point
│       └── main.go
├── internal/
│   ├── cmd/                 # Command implementations
│   │   ├── add.go
│   │   ├── delete.go
│   │   ├── list.go
│   │   ├── mark.go
│   │   └── root.go
│   ├── services/            # Business logic
│   │   └── json.go
│   └── utils/               # Utilities
│       └── log.go
└── README.md
```

## Testing

The application has comprehensive test coverage for all components:

- Command tests: Validates argument handling and command execution
- Service tests: Ensures task management functions work correctly
- Utility tests: Verifies logging functionality

Run the tests:

```
go test ./...
```

Generate coverage report:

```
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Configuration

Tasks are stored in a JSON file located at:
- Linux/macOS: `$HOME/.task-tracker.json`
- Windows: `%USERPROFILE%\.task-tracker.json`

## Releases

This project uses GitHub Actions to automatically build and publish releases. For more information, see:

- [RELEASES.md](RELEASES.md) - Detailed documentation on the release process
- [GitHub Releases Page](https://github.com/savabush/taskTracker/releases) - Download binaries

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 