package main

import (
	// "flag" // No longer needed if config manages log-file
	"fmt" // Package for formatted I/O (e.g., printing to console)
	"os"  // Package for operating system functionalities (e.g., exiting the program)

	// "strconv" // No longer needed in main.go
	// "strings" // No longer needed in main.go
	// Package for time-related functions (e.g., auto-save interval)
	"time"
)

const (
	// configPath is the default path for the application's configuration file.
	configPath = "config.json"
)

// main is the entry point of the CLI todo application.
// It initializes the logger, manages the todo list lifecycle (load, auto-save, save),
// and delegates command handling to the cli module.
func main() {
	// Load application configuration.
	config, err := LoadConfig(configPath)
	if err != nil {
		// If config loading fails, log the error and exit. No need to use SetupLogger yet,
		// as it might depend on the config itself. Just print to stderr.
		fmt.Fprintf(os.Stderr, "ERROR: Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	SetupLogger(config.LogFilePath) // Initialize the custom logger with potential log file from config.

	// Load the todo list from the data file specified in config.
	todoList, err := LoadFromFile(config.DataFile)
	if err != nil {
		// Log the error if loading fails and exit the application.
		LogError(err, "Failed to load todo list")
		PrintUserMessage("Error loading todo list. Exiting.")
		os.Exit(1) // Exit with an error code.
	}

	// Start a background goroutine for auto-saving the todo list periodically.
	// This ensures that changes are saved even if the application isn't explicitly exited.
	StartAutoSave(todoList, config.DataFile, time.Duration(config.AutoSaveInterval))

	// Delegate all command parsing and execution (both single command and interactive mode)
	// to the HandleCommands function in the cli module.
	HandleCommands(todoList)

	// Explicitly save the todo list to file before the application exits.
	// This is important for ensuring the latest changes are saved immediately,
	// especially for commands that don't trigger an auto-save shortly after.
	// This will also catch any changes made in interactive mode before the program fully terminates.
	err = todoList.SaveToFile(config.DataFile)
	if err != nil {
		// Log an error if saving fails during application shutdown.
		LogError(err, "Failed to save todo list on exit")
	}
}
