# Project: Git Branch Manager (gbm)
CLI tool that manages Git repository branches and worktrees based on configuration files, providing automated worktree management, JIRA integration, and intelligent file copying.

## Features
- **Automated Worktree Synchronization**: Manages worktrees based on `gbm.branchconfig.yaml` configuration
- **JIRA Integration**: Smart tab completion and branch naming from JIRA tickets
- **File Management**: Configurable file copying between worktrees with conflict handling
- **Shell Integration**: Tab completion and directory navigation (`gcd` function)
- **Worktree Lifecycle**: Create, list, remove, and switch between worktrees with safety checks
- **Repository Operations**: Bare repository support with branch validation and remote operations
- **Merge-back Management**: Intelligent merge-back suggestions and timestamp-based alerts
- **Configuration Validation**: Validates YAML syntax and branch references

## Tech Stack
- **Language**: Go 1.24.4
- **CLI Framework**: Cobra v1.9.1
- **Terminal UI**: Bubble Tea v1.3.5 with Lipgloss v1.1.0 styling
- **Git Operations**: go-git v5.16.2
- **Configuration**: TOML v1.5.0 and YAML v3.0.1
- **Testing**: testify v1.10.0 with custom Git test harness

## Structure
- `main.go` - Application entry point
- `cmd/` - CLI command implementations (13 commands)
- `internal/` - Core business logic (manager, git, config, state)
- `docs/` - Documentation and specifications
- `tools/` - Development utilities
- `todos/` - Task tracking

## Architecture
- **Layered Architecture**: Presentation (CLI), Business Logic, Data layers
- **Manager Pattern**: Central `Manager` struct coordinates all operations
- **Configuration Management**: Dual config system (YAML for worktrees, TOML for settings)
- **Git Abstraction**: `GitManager` provides clean interface over Git operations

## Commands
- Build: `go build -o gbm .`
- Test: `go test ./...`
- Lint: `golangci-lint run`
- Dev/Run: `go run .`

## Testing
- **Framework**: Uses testify for assertions and test suites
- **Git Test Harness**: Sophisticated testing environment with pre-configured scenarios
- **Test Helpers**: Shared utilities in `cmd/test_helpers.go`
- **Coverage**: Run `go test -coverprofile=coverage.out ./...` then `go tool cover -html=coverage.out`