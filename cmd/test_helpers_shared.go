package cmd

import (
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
