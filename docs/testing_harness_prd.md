# Product Requirements Document: Git Testing Harness Infrastructure

## Overview

A testing harness infrastructure for the Git Branch Manager (gbm) CLI tool that creates isolated, local Git repositories for testing. This document focuses solely on the harness setup and infrastructure - not the actual test implementations. The harness will provide the foundation for running gbm commands against realistic Git environments without external dependencies.

## Objectives

### Primary Goals
- **Isolated Test Environments**: Create completely isolated Git repositories that don't affect external systems
- **Realistic Git Scenarios**: Simulate real-world Git workflows with branches, remotes, and various repository states
- **Fast Setup**: Enable rapid test environment creation and teardown
- **Zero External Dependencies**: No network or external service requirements for test execution
- **Flexible Configuration**: Support various repository configurations and scenarios

### Success Metrics
- Test environment setup time < 100ms per repository
- Support for all Git operations needed by gbm
- Zero external dependencies for test execution
- Clean teardown with no leftover temporary files

## Technical Requirements

### Core Dependencies
- **testify/assert**: For fluent assertions in future tests
- **testify/require**: For test requirements that should halt execution on failure
- **Go testing package**: Standard testing framework integration
- **os/exec**: For Git command execution
- **path/filepath**: For cross-platform path handling

### Architecture Components

#### 1. GitTestRepo Structure
```go
type GitTestRepo struct {
    RemoteDir    string           // Path to bare remote repository
    LocalDir     string           // Path to local working repository
    TempDir      string           // Root temporary directory
    Config       RepoConfig       // Repository configuration
    t           *testing.T       // Test context
}

type RepoConfig struct {
    DefaultBranch    string
    UserName         string
    UserEmail        string
    RemoteName       string
}
```

#### 2. Repository Setup Methods

**Core Setup**
- `NewGitTestRepo(t *testing.T, opts ...RepoOption) *GitTestRepo`
- `setupBareRemote() error`
- `setupLocalRepo() error`
- `configureGitUser() error`
- `createInitialCommit() error`
- `Cleanup()`

**Configuration Options**
```go
type RepoOption func(*GitTestRepo)

func WithDefaultBranch(branch string) RepoOption
func WithUser(name, email string) RepoOption
func WithRemoteName(name string) RepoOption
```

#### 3. Content Creation Methods

**Branch Operations**
- `CreateBranch(name, content string) error`
- `CreateBranchFrom(name, baseBranch, content string) error`
- `SwitchToBranch(name string) error`

**File Operations**
- `WriteFile(path, content string) error`
- `CreateEnvrc(mapping map[string]string) error`
- `CommitChanges(message string) error`
- `PushBranch(branch string) error`

**Repository State**
- `ConvertToBare() string`
- `GetLocalPath() string`
- `GetRemotePath() string`
- `ListBranches() ([]string, error)`

#### 4. Mock Service Infrastructure

**JIRA CLI Mock**
```go
type MockJiraCLI struct {
    responses map[string]string
    commands  []string
}

func NewMockJiraCLI() *MockJiraCLI
func (m *MockJiraCLI) SetResponse(command, response string)
func (m *MockJiraCLI) InstallMock() error
func (m *MockJiraCLI) RemoveMock() error
```

**Environment Mocking**
- Mock `jira` CLI command availability
- Simulate command execution environments
- Control external command responses

## Implementation Specification

### File Structure
```
internal/
├── testutils/
│   ├── git_harness.go          # Main GitTestRepo implementation
│   ├── repo_options.go         # Configuration options
│   ├── mock_services.go        # Mock external services
│   └── scenarios.go            # Pre-defined test scenarios
```

### Core Implementation Requirements

#### 1. GitTestRepo Core Methods

**Constructor with Options**
```go
func NewGitTestRepo(t *testing.T, opts ...RepoOption) *GitTestRepo {
    // Create temp directory using t.TempDir()
    // Apply configuration options
    // Setup bare remote repository
    // Clone to create local repository
    // Configure git user and email
    // Create initial commit and push
    // Return configured GitTestRepo
}
```

**Cleanup Implementation**
```go
func (r *GitTestRepo) Cleanup() {
    // Cleanup handled automatically by t.TempDir()
    // Remove any global state changes
    // Restore environment variables
}
```

#### 2. Repository Setup Flow

**Bare Remote Setup**
1. Create temporary directory for bare repository
2. Initialize bare Git repository (`git init --bare`)
3. Configure remote repository settings

**Local Repository Setup**
1. Clone bare repository to local directory
2. Configure git user name and email
3. Create initial README.md file
4. Make initial commit
5. Push to remote repository

**Branch Creation Workflow**
1. Switch to appropriate base branch
2. Create new branch (`git checkout -b`)
3. Create/modify files with provided content
4. Stage and commit changes
5. Push branch to remote repository
6. Return to original branch

#### 3. Configuration Management

**Default Configuration**
```go
var defaultConfig = RepoConfig{
    DefaultBranch: "main",
    UserName:      "Test User",
    UserEmail:     "test@example.com",
    RemoteName:    "origin",
}
```

**Option Pattern Implementation**
```go
func WithDefaultBranch(branch string) RepoOption {
    return func(r *GitTestRepo) {
        r.Config.DefaultBranch = branch
    }
}

func WithUser(name, email string) RepoOption {
    return func(r *GitTestRepo) {
        r.Config.UserName = name
        r.Config.UserEmail = email
    }
}
```

#### 4. Error Handling Strategy

**Git Command Execution**
- Wrap all git commands with proper error handling
- Capture both stdout and stderr for debugging
- Provide context about which operation failed
- Include command and arguments in error messages

**Path Management**
- Use filepath.Join for cross-platform compatibility
- Validate paths before operations
- Handle long path limitations on Windows
- Ensure proper cleanup of temporary directories

### Mock Service Requirements

#### 1. JIRA CLI Simulation

**Command Response Mapping**
```go
type JiraResponse struct {
    Command  string
    Output   string
    ExitCode int
}

func (m *MockJiraCLI) AddResponse(command, output string, exitCode int)
func (m *MockJiraCLI) SimulateFailure(command string, errorMsg string)
```

**Installation Mechanism**
- Create temporary executable that responds to `jira` commands
- Modify PATH to include mock executable directory
- Restore original PATH on cleanup

#### 2. Environment Control

**Command Availability Simulation**
- Simulate presence/absence of external commands
- Control command execution responses
- Mock network connectivity issues
- Simulate permission errors

## Pre-defined Scenarios

### Standard Repository Configurations

#### 1. Basic Repository
```go
func NewBasicRepo(t *testing.T) *GitTestRepo {
    return NewGitTestRepo(t,
        WithDefaultBranch("main"),
        WithUser("Test User", "test@example.com"),
    )
}
```

#### 2. Multi-Branch Repository
```go
func NewMultiBranchRepo(t *testing.T) *GitTestRepo {
    repo := NewBasicRepo(t)
    repo.CreateBranch("develop", "Development content")
    repo.CreateBranch("feature/auth", "Authentication feature")
    repo.CreateBranch("production/v1.0", "Production release")
    return repo
}
```

#### 3. Repository with .envrc
```go
func NewEnvrcRepo(t *testing.T, mapping map[string]string) *GitTestRepo {
    repo := NewMultiBranchRepo(t)
    repo.CreateEnvrc(mapping)
    repo.CommitChanges("Add .envrc configuration")
    repo.PushBranch("main")
    return repo
}
```

### Utility Functions

#### 1. Directory Navigation
```go
func (r *GitTestRepo) WithWorkingDir(dir string) func() {
    // Change to specified directory
    // Return function to restore original directory
}

func (r *GitTestRepo) InLocalRepo(fn func() error) error {
    // Execute function with LocalDir as working directory
    // Restore original working directory after execution
}
```

#### 2. Command Execution Helpers
```go
func (r *GitTestRepo) runGitCommand(args ...string) error {
    // Execute git command in LocalDir
    // Handle errors with context
    // Log command execution for debugging
}

func (r *GitTestRepo) runCommand(name string, args ...string) ([]byte, error) {
    // Execute arbitrary command in LocalDir
    // Return output and error
}
```

## Usage Examples

### Basic Usage
```go
func TestSomeFeature(t *testing.T) {
    // Setup test repository
    repo := testutils.NewGitTestRepo(t,
        testutils.WithDefaultBranch("main"),
        testutils.WithUser("Test User", "test@example.com"),
    )

    // Create test branches
    repo.CreateBranch("develop", "Development branch content")
    repo.CreateBranch("feature/new-feature", "New feature content")

    // Create .envrc configuration
    repo.CreateEnvrc(map[string]string{
        "MAIN": "main",
        "DEV":  "develop",
        "FEAT": "feature/new-feature",
    })

    // Test gbm commands in repository
    originalDir, _ := os.Getwd()
    defer os.Chdir(originalDir)
    os.Chdir(repo.GetLocalPath())

    // Run gbm commands and verify results
    // (actual test implementation not included in this PRD)
}
```

### Clone Testing Setup
```go
func TestCloneCommand(t *testing.T) {
    // Setup source repository
    repo := testutils.NewMultiBranchRepo(t)
    repo.CreateEnvrc(map[string]string{
        "MAIN": "main",
        "DEV":  "develop",
    })

    // Convert to bare repository for cloning
    bareRepoPath := repo.ConvertToBare()

    // Create target directory for clone
    targetDir := filepath.Join(t.TempDir(), "cloned-repo")

    // Test gbm clone command
    // (actual test implementation not included in this PRD)
}
```

### Mock JIRA Integration
```go
func TestJiraIntegration(t *testing.T) {
    // Setup mock JIRA CLI
    mockJira := testutils.NewMockJiraCLI()
    defer mockJira.RemoveMock()

    // Configure mock responses
    mockJira.SetResponse("jira me", "test-user@example.com")
    mockJira.SetResponse("jira issue list -atest-user@example.com --plain",
        "Story\tPROJ-123\tUser Authentication\tIn Progress")

    mockJira.InstallMock()

    // Setup repository
    repo := testutils.NewGitTestRepo(t)

    // Test JIRA-integrated gbm commands
    // (actual test implementation not included in this PRD)
}
```

## Acceptance Criteria

### Core Functionality
- ✅ Create isolated Git repositories with remote/local setup
- ✅ Support branch creation, file operations, and git commands
- ✅ Provide .envrc creation and management utilities
- ✅ Clean up all temporary resources automatically
- ✅ Cross-platform compatibility (Windows, macOS, Linux)

### Performance Requirements
- ✅ Repository setup time < 100ms
- ✅ Clean teardown with no leftover files
- ✅ Memory usage < 10MB per test repository
- ✅ Support for concurrent test execution

### Mock Services
- ✅ JIRA CLI command mocking capability
- ✅ Configurable command responses and failures
- ✅ Environment variable control
- ✅ Proper cleanup of mock installations

### Developer Experience
- ✅ Simple, intuitive API for repository creation
- ✅ Comprehensive error messages for setup failures
- ✅ Clear documentation and usage examples
- ✅ Minimal boilerplate for basic scenarios

### Quality Standards
- ✅ All public methods documented with examples
- ✅ Proper error handling for all operations
- ✅ Cross-platform path handling
- ✅ Integration with Go testing framework and testify

This testing harness infrastructure will provide the foundation for comprehensive testing of gbm without requiring external dependencies or complex setup procedures.
