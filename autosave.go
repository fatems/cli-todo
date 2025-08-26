package main

import (
	"time" // Package for time-related operations, used for `time.After` and `time.Duration`
)

// StartAutoSave Goroutine initiates a background process that periodically saves
// the current state of the TodoList to a specified JSON file.
// It takes a pointer to the TodoList, the filename for persistence, and the interval
// at which to perform the auto-save operation.
func StartAutoSave(todoList *TodoList, filename string, interval time.Duration) {
	// The `go func()` syntax starts a new goroutine, allowing the auto-save logic
	// to run concurrently with the main application flow without blocking it.
	go func() {
		// This infinite loop ensures the auto-save runs continuously until the application exits.
		for {
			// `<-time.After(interval)` blocks the goroutine for the specified `interval`.
			// After the interval, it sends a value to the channel returned by `time.After`,
			// unblocking the goroutine and allowing the next iteration to proceed.
			<-time.After(interval)

			// Attempt to save the TodoList to the file.
			err := todoList.SaveToFile(filename)
			if err != nil {
				// If saving fails, log an error with a descriptive message.
				LogError(err, "Auto-save failed")
			} else {
				// If saving is successful, log an informational message.
				LogInfo("Auto-saved todo list.")
			}
		}
	}()
}
