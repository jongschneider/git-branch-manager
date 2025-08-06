package cmd

import (
	"bytes"
	"io"
	"os"
)

// Helper function to simulate user input for confirmation prompts
func simulateUserInput(input string, fn func() error) error {
	// Create a pipe to simulate stdin
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r

	// Write the input
	go func() {
		defer func() { _ = w.Close() }()
		_, _ = w.Write([]byte(input + "\n"))
	}()

	// Execute the function
	err := fn()

	// Restore stdin
	os.Stdin = oldStdin
	_ = r.Close()

	return err
}

// Helper function to capture command output
func captureOutput(fn func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String(), err
}
