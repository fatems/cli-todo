package main

import (
	"bufio"   // Package for buffered I/O operations (e.g., reading from stdin)
	"flag"    // Package for parsing command-line flags
	"fmt"     // Package for formatted I/O (e.g., printing to console)
	"os"      // Package for operating system functionalities (e.g., exiting the program)
	"strconv" // Package for converting strings to other data types
	"strings" // Package for string manipulation
	"time"    // Package for handling dates and times
)

// ActionType represents the type of action performed.
type ActionType int

const (
	ActionNone ActionType = iota
	ActionAdd
	ActionComplete
	ActionDelete
	ActionUncomplete
)

// lastAction stores information about the last performed action for undo functionality.
type lastAction struct {
	Type ActionType
	ID   int // ID of the todo affected by the action
	// For delete, we need to store the entire Todo object to re-add it.
	DeletedTodo *Todo
	// For complete/uncomplete, we need to store the previous completed status.
	PreviousCompletedStatus bool
}

// lastActionState tracks the most recent action for undo purposes.
var lastActionState lastAction

// runInteractiveMode provides a continuous loop for user interaction,
// prompting for commands and executing them until the user decides to exit.
// It directly interacts with the TodoList and utilizes logging utilities.
func runInteractiveMode(todoList *TodoList) {
	PrintUserMessage("🚀 Entering interactive mode. Type 'help' for commands, 'exit' to quit.")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")                     // Keep prompt on stdout
		input, _ := reader.ReadString('\n') // Read user input until a newline character.
		command := strings.TrimSpace(input) // Remove leading/trailing whitespace.

		splitCommand := strings.Fields(command) // Split the command string into fields.
		if len(splitCommand) == 0 {
			continue // If input is empty, prompt again.
		}

		subCommand := strings.ToLower(splitCommand[0]) // Get the main command (e.g., "add", "list").

		switch subCommand {
		case "add":
			// Interactive add command needs to parse task, priority, due date, and tags from the input string.
			// This is a simplified approach; a dedicated interactive parser would be more robust.
			parts := splitCommand[1:]
			task := ""
			priority := ""
			dueDateStr := ""
			tags := []string{}

			for i := 0; i < len(parts); i++ {
				if parts[i] == "-p" && i+1 < len(parts) {
					priority = parts[i+1]
					i++
				} else if parts[i] == "-d" && i+1 < len(parts) {
					dueDateStr = parts[i+1]
					i++
				} else if parts[i] == "-t" && i+1 < len(parts) {
					tags = append(tags, strings.Split(parts[i+1], ",")...)
					i++
				} else {
					if task == "" { // First unflagged part is the task
						task = parts[i]
					} else {
						task += " " + parts[i]
					}
				}
			}

			if task == "" {
				PrintUserMessage("Usage: add <task> [-p <priority>] [-d <YYYY-MM-DD>] [-t <tag1,tag2>]")
				LogError(fmt.Errorf("missing task for add command"), "Interactive mode input error")
			} else {
				var dueDate *time.Time
				if dueDateStr != "" {
					parsedDate, err := parseDueDate(dueDateStr)
					if err != nil {
						PrintUserMessage("Invalid due date format. Use YYYY-MM-DD.")
						LogError(err, "Interactive mode input error: invalid due date")
						continue
					}
					dueDate = &parsedDate
				}
				todoList.Add(task, toCanonicalPriority(PriorityLevel(priority)), dueDate, tags)
				lastActionState = lastAction{Type: ActionAdd, ID: todoList.NextID - 1} // Store ID of newly added todo
			}
		case "edit":
			if len(splitCommand) < 3 {
				PrintUserMessage("Usage: edit <id> <new_task_description>")
				LogError(fmt.Errorf("missing ID or new task for edit command"), "Interactive mode input error")
			} else {
				id, err := strconv.Atoi(splitCommand[1])
				if err != nil {
					PrintUserMessage("Invalid ID. Please provide a number.")
					LogError(err, "Interactive mode input error: invalid ID for edit")
				} else {
					newTask := strings.Join(splitCommand[2:], " ")
					err = todoList.EditTask(id, newTask)
					if err != nil {
						LogError(err, fmt.Sprintf("Failed to edit todo with ID %d in interactive mode", id))
						PrintUserMessage(err.Error())
					}
				}
			}
		case "clear-completed":
			if getConfirmation("Are you sure you want to clear all completed todos?") {
				todoList.ClearCompleted()
			} else {
				PrintUserMessage("Clearing completed todos cancelled.")
			}
		case "search":
			if len(splitCommand) < 2 {
				PrintUserMessage("Usage: search <query>")
				LogError(fmt.Errorf("missing query for search command"), "Interactive mode input error")
			} else {
				query := strings.Join(splitCommand[1:], " ")
				results := todoList.SearchTasks(query)
				if len(results.Todos) == 0 {
					PrintUserMessage(fmt.Sprintf("🔍 No tasks found matching \"%s\".", query))
				} else {
					PrintUserMessage(fmt.Sprintf("🔍 Tasks matching \"%s\":", query))
					results.List(ListOptions{}) // List with default options for search results
				}
			}
		case "complete":
			if len(splitCommand) < 2 {
				PrintUserMessage("Usage: complete <id>")
				LogError(fmt.Errorf("missing ID for complete command"), "Interactive mode input error")
			} else {
				id, err := strconv.Atoi(splitCommand[1])
				if err != nil {
					PrintUserMessage("Invalid ID. Please provide a number.")
					LogError(err, "Interactive mode input error: invalid ID for complete")
				} else {
					err = todoList.Complete(id)
					if err != nil {
						LogError(err, fmt.Sprintf("Failed to complete todo with ID %d in interactive mode", id))
						PrintUserMessage(err.Error())
					} else {
						// Assuming completed status was false before completing.
						lastActionState = lastAction{Type: ActionComplete, ID: id, PreviousCompletedStatus: false}
					}
				}
			}
		case "uncomplete": // New command for undo functionality
			if len(splitCommand) < 2 {
				PrintUserMessage("Usage: uncomplete <id>")
				LogError(fmt.Errorf("missing ID for uncomplete command"), "Interactive mode input error")
			} else {
				id, err := strconv.Atoi(splitCommand[1])
				if err != nil {
					PrintUserMessage("Invalid ID. Please provide a number.")
					LogError(err, "Interactive mode input error: invalid ID for uncomplete")
				} else {
					err = todoList.Uncomplete(id)
					if err != nil {
						LogError(err, fmt.Sprintf("Failed to uncomplete todo with ID %d in interactive mode", id))
						PrintUserMessage(err.Error())
					} else {
						// Assuming completed status was true before uncompleting.
						lastActionState = lastAction{Type: ActionUncomplete, ID: id, PreviousCompletedStatus: true}
					}
				}
			}
		case "delete":
			if len(splitCommand) < 2 {
				PrintUserMessage("Usage: delete <id>")
				LogError(fmt.Errorf("missing ID for delete command"), "Interactive mode input error")
			} else {
				id, err := strconv.Atoi(splitCommand[1])
				if err != nil {
					PrintUserMessage("Invalid ID. Please provide a number.")
					LogError(err, "Interactive mode input error: invalid ID for delete")
				} else {
					if getConfirmation("Are you sure you want to delete todo with ID " + strconv.Itoa(id) + "?") {
						deletedTodo, err := todoList.Delete(id)
						if err != nil {
							LogError(err, fmt.Sprintf("Failed to delete todo with ID %d in interactive mode", id))
							PrintUserMessage(err.Error())
						} else {
							lastActionState = lastAction{Type: ActionDelete, ID: id, DeletedTodo: &deletedTodo}
						}
					} else {
						PrintUserMessage(fmt.Sprintf("Deletion of todo #%d cancelled.", id))
					}
				}
			}
		case "list":
			// For enhanced list, we'll need to parse additional flags here in interactive mode
			// For now, just call simple list.
			todoList.List(ListOptions{}) // Display all current todos with default options for now.
		case "help":
			// Print available commands for interactive mode.
			PrintUserMessage("✨ Commands:")
			PrintUserMessage("  ➕ add <task> [-p <high|medium|low>] [-d <YYYY-MM-DD>] [-t <tag1,tag2>]  - Add a new todo task")
			PrintUserMessage("  ✏️ edit <id> <new_task>                                            - Edit a task description")
			PrintUserMessage("  🧹 clear-completed                                                 - Remove all completed todos")
			PrintUserMessage("  🔍 search <query>                                                  - Search tasks by description or tags")
			PrintUserMessage("  🔄 uncomplete <id>                                                - Mark a todo as incomplete by ID")
			PrintUserMessage("  ↩️ undo                                                             - Undo the last action")
			PrintUserMessage("  ✅ complete <id>                                                  - Mark a todo as complete by ID")
			PrintUserMessage("  🗑️ delete <id>                                                    - Delete a todo by ID")
			PrintUserMessage("  📋 list                                                           - List all todos")
			PrintUserMessage("  🚪 exit                                                           - Exit interactive mode")
		case "exit":
			PrintUserMessage("👋 Exiting interactive mode.")
			return // Exit the interactive loop.
		case "undo": // New undo command
			switch lastActionState.Type {
			case ActionAdd:
				deletedTodo, err := todoList.Delete(lastActionState.ID) // Undo add is a delete
				if err != nil {
					LogError(err, fmt.Sprintf("Failed to undo add for todo ID %d", lastActionState.ID))
					PrintUserMessage("❌ Undo failed: " + err.Error())
				} else {
					PrintUserMessage(fmt.Sprintf("↩️ Undid adding todo #%d (task: \"%s\").", lastActionState.ID, deletedTodo.Task))
				}
			case ActionComplete:
				err := todoList.Uncomplete(lastActionState.ID)
				if err != nil {
					LogError(err, fmt.Sprintf("Failed to undo complete for todo ID %d", lastActionState.ID))
					PrintUserMessage("❌ Undo failed: " + err.Error())
				} else {
					PrintUserMessage(fmt.Sprintf("↩️ Undid completing todo #%d.", lastActionState.ID))
				}
			case ActionDelete:
				if lastActionState.DeletedTodo != nil {
					// To undo delete, we re-add the todo with its original state.
					// Note: This will assign a *new* ID if NextID has advanced. For true undo, we'd need to re-insert at original ID.
					// For basic undo, re-adding is sufficient.
					todoList.Todos = append(todoList.Todos, *lastActionState.DeletedTodo)
					PrintUserMessage(fmt.Sprintf("↩️ Undid deleting todo #%d (re-added as #%d: \"%s\").", lastActionState.ID, lastActionState.DeletedTodo.ID, lastActionState.DeletedTodo.Task))
				} else {
					PrintUserMessage("❌ Cannot undo delete: no todo data stored.")
					LogError(fmt.Errorf("attempted to undo delete without stored todo data"), "Undo error")
				}
			case ActionUncomplete:
				err := todoList.Complete(lastActionState.ID)
				if err != nil {
					LogError(err, fmt.Sprintf("Failed to undo uncomplete for todo ID %d", lastActionState.ID))
					PrintUserMessage("❌ Undo failed: " + err.Error())
				} else {
					PrintUserMessage(fmt.Sprintf("↩️ Undid uncompleting todo #%d.", lastActionState.ID))
				}
			case ActionNone:
				PrintUserMessage("🤔 No action to undo.")
			}
			lastActionState.Type = ActionNone // Clear the last action after undo
			// Note: clearing ID and DeletedTodo might also be good here depending on desired robustness.
			lastActionState.ID = 0
			lastActionState.DeletedTodo = nil
			lastActionState.PreviousCompletedStatus = false
		default:
			PrintUserMessage("❓ Unknown command. Type 'help' for a list of commands.")
			LogError(fmt.Errorf("unknown command: %s", subCommand), "Interactive mode input error")
		}
	}
}

// parseDueDate parses a date string in YYYY-MM-DD format into a time.Time object.
func parseDueDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// getConfirmation prompts the user for a yes/no confirmation and returns true if 'y' or 'Y' is entered.
func getConfirmation(prompt string) bool {
	fmt.Printf("%s (y/N): ", prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}

// processSingleCommand handles the execution of a single command based on the provided flags.
// It takes the TodoList and pointers to the parsed flag values.
func processSingleCommand(todoList *TodoList, addPtr *string, completePtr *int, deletePtr *int, listPtr *bool, clearCompletedCmdPtr *bool, filterStatusPtr *string, filterPriorityPtr *string, filterTagsPtr *string, sortByPtr *string, sortOrderPtr *string) {
	switch {
	case flag.NFlag() == 0:
		// If no flags are provided at all, print usage and suggest interactive mode.
		PrintUserMessage("Usage: go run . [options]")
		PrintUserMessage("💡 Run with -interactive for interactive mode.")
		flag.PrintDefaults() // Display default values and descriptions for all flags.
	case *addPtr != "":
		// If the -add flag is present, add a new todo with the provided task description.
		// For single command mode, priority, due date, and tags are not yet supported via flags directly.
		todoList.Add(*addPtr, toCanonicalPriority(PriorityMedium), nil, []string{}) // Default values for new fields
	case *completePtr != 0:
		// If the -complete flag is present, mark the todo with the given ID as complete.
		err := todoList.Complete(*completePtr)
		if err != nil {
			// Log and print an error if the todo to complete is not found.
			LogError(err, fmt.Sprintf("Failed to complete todo with ID %d", *completePtr))
			PrintUserMessage(err.Error())
		}
	case *deletePtr != 0:
		// If the -delete flag is present, remove the todo with the given ID.
		if getConfirmation(fmt.Sprintf("Are you sure you want to delete todo with ID %d?", *deletePtr)) {
			deletedTodo, err := todoList.Delete(*deletePtr)
			if err != nil {
				// Log and print an error if the todo to delete is not found.
				LogError(err, fmt.Sprintf("Failed to delete todo with ID %d", *deletePtr))
				PrintUserMessage(err.Error())
			} else {
				// No undo state stored for single commands for simplicity here.
				_ = deletedTodo // Use deletedTodo to avoid unused variable error
			}
		} else {
			PrintUserMessage(fmt.Sprintf("Deletion of todo #%d cancelled.", *deletePtr))
		}
	case *clearCompletedCmdPtr:
		// If the -clear-completed flag is present, clear all completed todos.
		if getConfirmation("Are you sure you want to clear all completed todos?") {
			todoList.ClearCompleted()
		} else {
			PrintUserMessage("Clearing completed todos cancelled.")
		}
	case *listPtr:
		// If the -list flag is present, display all current todos with applied filters and sorting.
		options := ListOptions{
			FilterStatus:   *filterStatusPtr,
			FilterPriority: PriorityLevel(*filterPriorityPtr),
			FilterTags:     strings.Split(*filterTagsPtr, ","),
			SortBy:         *sortByPtr,
			SortOrder:      *sortOrderPtr,
		}
		// Clean up empty tag strings from splitting
		if len(options.FilterTags) == 1 && options.FilterTags[0] == "" {
			options.FilterTags = []string{}
		}
		todoList.List(options)
	default:
		// This case catches any other combination of flags that don't match specific commands.
		PrintUserMessage("❌ Unknown command or invalid flag combination. Type 'go run .' for usage.")
	}
}

// HandleCommands parses command-line flags and manages the application flow,
// either by executing a single command or entering an interactive mode.
// It defines the CLI flags, parses them, and then dispatches control
// to either `runInteractiveMode` or `processSingleCommand` based on user input.
func HandleCommands(todoList *TodoList) {
	// Define command-line flags for various todo operations.
	addPtr := flag.String("add", "", "Add a new todo task")
	completePtr := flag.Int("complete", 0, "Mark a todo as complete by ID")
	deletePtr := flag.Int("delete", 0, "Delete a todo by ID")
	listPtr := flag.Bool("list", false, "List all todos")
	interactivePtr := flag.Bool("interactive", false, "Run in interactive mode")
	clearCompletedCmdPtr := flag.Bool("clear-completed", false, "Clear all completed todos") // New flag for single command

	// New flags for enhanced list command
	filterStatusPtr := flag.String("filter-status", "all", "Filter todos by status (all, completed, incomplete)")
	filterPriorityPtr := flag.String("filter-priority", "", "Filter todos by priority (high, medium, low)")
	filterTagsPtr := flag.String("filter-tags", "", "Filter todos by tags (comma-separated, e.g., work,urgent)")
	sortByPtr := flag.String("sort-by", "id", "Sort todos by field (id, task, created_at, due_date, priority)")
	sortOrderPtr := flag.String("sort-order", "asc", "Sort order (asc, desc)")

	flag.Parse() // Parse the command-line arguments into the defined flags.

	// If interactive mode is enabled, run the interactive loop.
	if *interactivePtr {
		runInteractiveMode(todoList)
		return // Exit after interactive mode finishes
	}

	// If not in interactive mode, process a single command based on the provided flags.
	processSingleCommand(todoList, addPtr, completePtr, deletePtr, listPtr, clearCompletedCmdPtr, filterStatusPtr, filterPriorityPtr, filterTagsPtr, sortByPtr, sortOrderPtr)
}
