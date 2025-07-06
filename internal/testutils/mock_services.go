package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type MockJiraCLI struct {
	responses  map[string]JiraResponse
	mockDir    string
	originalPATH string
	t          *testing.T
}

type JiraResponse struct {
	Command  string
	Output   string
	ExitCode int
}

func NewMockJiraCLI(t *testing.T) *MockJiraCLI {
	return &MockJiraCLI{
		responses: make(map[string]JiraResponse),
		t:         t,
	}
}

func (m *MockJiraCLI) SetResponse(command, output string) {
	m.responses[command] = JiraResponse{
		Command:  command,
		Output:   output,
		ExitCode: 0,
	}
}

func (m *MockJiraCLI) AddResponse(command, output string, exitCode int) {
	m.responses[command] = JiraResponse{
		Command:  command,
		Output:   output,
		ExitCode: exitCode,
	}
}

func (m *MockJiraCLI) SimulateFailure(command string, errorMsg string) {
	m.responses[command] = JiraResponse{
		Command:  command,
		Output:   errorMsg,
		ExitCode: 1,
	}
}

func (m *MockJiraCLI) InstallMock() error {
	tempDir, err := os.MkdirTemp("", "jira-mock-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	m.mockDir = tempDir
	m.originalPATH = os.Getenv("PATH")

	mockScript := m.generateMockScript()

	var mockPath string
	if runtime.GOOS == "windows" {
		mockPath = filepath.Join(tempDir, "jira.bat")
	} else {
		mockPath = filepath.Join(tempDir, "jira")
	}

	if err := os.WriteFile(mockPath, []byte(mockScript), 0755); err != nil {
		return fmt.Errorf("failed to write mock script: %w", err)
	}

	newPATH := tempDir + string(os.PathListSeparator) + m.originalPATH
	m.t.Setenv("PATH", newPATH)

	return nil
}

func (m *MockJiraCLI) RemoveMock() error {
	if m.mockDir != "" {
		if err := os.RemoveAll(m.mockDir); err != nil {
			return fmt.Errorf("failed to remove mock directory: %w", err)
		}
	}

	return nil
}

func (m *MockJiraCLI) generateMockScript() string {
	if runtime.GOOS == "windows" {
		return m.generateWindowsScript()
	}
	return m.generateUnixScript()
}

func (m *MockJiraCLI) generateUnixScript() string {
	script := "#!/bin/bash\n\n"
	script += "COMMAND=\"$*\"\n\n"

	for _, response := range m.responses {
		script += fmt.Sprintf("if [ \"$COMMAND\" = \"%s\" ]; then\n", response.Command)
		script += fmt.Sprintf("    echo '%s'\n", strings.ReplaceAll(response.Output, "'", "'\"'\"'"))
		script += fmt.Sprintf("    exit %d\n", response.ExitCode)
		script += "fi\n\n"
	}

	script += "echo \"Mock JIRA CLI: Unknown command: $COMMAND\" >&2\n"
	script += "exit 1\n"

	return script
}

func (m *MockJiraCLI) generateWindowsScript() string {
	script := "@echo off\n\n"
	script += "set COMMAND=%*\n\n"

	for _, response := range m.responses {
		script += fmt.Sprintf("if \"%%COMMAND%%\" == \"%s\" (\n", response.Command)
		script += fmt.Sprintf("    echo %s\n", response.Output)
		script += fmt.Sprintf("    exit /b %d\n", response.ExitCode)
		script += ")\n\n"
	}

	script += "echo Mock JIRA CLI: Unknown command: %COMMAND% >&2\n"
	script += "exit /b 1\n"

	return script
}

func IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func MockCommandUnavailable(t *testing.T, command string) {
	originalPATH := os.Getenv("PATH")

	paths := strings.Split(originalPATH, string(os.PathListSeparator))
	var filteredPaths []string

	for _, path := range paths {
		if !strings.Contains(path, command) {
			filteredPaths = append(filteredPaths, path)
		}
	}

	t.Setenv("PATH", strings.Join(filteredPaths, string(os.PathListSeparator)))
}