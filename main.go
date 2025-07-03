package main

import (
	"os"

	"gbm/cmd"
)

func main() {
	defer cmd.CloseLogFile()

	if err := cmd.Execute(); err != nil {
		cmd.PrintError("Error: %v", err)
		os.Exit(1)
	}
}

