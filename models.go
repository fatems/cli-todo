package main

import (
	"encoding/json" // Package for JSON encoding and decoding
	"fmt"           // Package for formatted I/O (e.g., error messages, print statements)
	"os"            // Package for operating system functionalities (e.g., file operations)
	"sort"          // Package for sorting slices
	"strings"       // Package for string manipulation (e.g., Contains, ToLower)
	"time"          // Package for time-related functions (e.g., todo creation timestamp)
)

// PriorityLevel defines the type for todo priority, ensuring valid values.
type PriorityLevel string

// Constants for predefined priority levels.
const (
	PriorityHigh   PriorityLevel = "high"
	PriorityMedium PriorityLevel = "medium"
	PriorityLow    PriorityLevel = "low"
)

// priorityMap provides a case-insensitive lookup for canonical PriorityLevel values.
var priorityMap = map[string]PriorityLevel{
	"high":   PriorityHigh,
	"medium": PriorityMedium,
	"low":    PriorityLow,
}

// Todo represents a single task item in the todo list.
// It includes fields for a unique identifier, the task description, its completion status,
// and the timestamp of its creation.
type Todo struct {
	ID        int           `json:"id"`         // Unique identifier for the todo item.
	Task      string        `json:"task"`       // The description of the task.
	Completed bool          `json:"completed"`  // A boolean indicating if the task is completed (true) or not (false).
	CreatedAt time.Time     `json:"created_at"` // The timestamp when the todo item was created.
	Priority  PriorityLevel `json:"priority"`   // Priority of the todo (e.g., "high", "medium", "low").
	DueDate   *time.Time    `json:"due_date"`   // Optional due date for the todo item.
	Tags      []string      `json:"tags"`       // Optional tags/categories for the todo item.
}

// TodoList manages a collection of Todo items.
// It holds a slice of `Todo` structs and keeps track of the next available ID
// to ensure uniqueness for new todo items.
type TodoList struct {
	Todos  []Todo `json:"todos"`   // A slice (dynamic array) of Todo items.
	NextID int    `json:"next_id"` // The next ID to be assigned to a new todo item. This ensures unique IDs.
}

// NewTodoList creates and returns a pointer to a new, empty TodoList.
// The list is initialized with an empty slice of todos and a starting NextID of 1.
func NewTodoList() *TodoList {
	return &TodoList{
		Todos:  []Todo{}, // Initialize with an empty slice.
		NextID: 1,        // Start IDs from 1.
	}
}

// Add a new todo item to the TodoList.
// It takes a task description as input, creates a new Todo struct with a unique ID,
// sets its status to incomplete, records the creation time, and appends it to the list.
func (tl *TodoList) Add(task string, priority PriorityLevel, dueDate *time.Time, tags []string) {
	// Normalize the priority input to a canonical form.
	canonicalPriority := toCanonicalPriority(priority)

	// Validate canonical priority level. If invalid, default to medium and log a warning.
	if !isValidPriority(canonicalPriority) {
		LogWarning(fmt.Sprintf("Invalid priority level '%s' for task '%s'. Defaulting to Medium.", priority, task))
		canonicalPriority = PriorityMedium
	}

	// Create a new Todo instance.
	todo := Todo{
		ID:        tl.NextID,  // Assign the next available ID.
		Task:      task,       // Set the provided task description.
		Completed: false,      // New tasks are incomplete by default.
		CreatedAt: time.Now(), // Record the current time.
		Priority:  canonicalPriority,
		DueDate:   dueDate,
		Tags:      tags,
	}
	// Append the new todo to the existing slice of todos.
	tl.Todos = append(tl.Todos, todo)
	tl.NextID++ // Increment NextID for the next new todo.
	PrintUserMessage(fmt.Sprintf("✅ Added todo #%d: \"%s\"", todo.ID, todo.Task))
}

// isValidPriority checks if the given priority level is one of the predefined valid levels.
func isValidPriority(p PriorityLevel) bool {
	switch p {
	case PriorityHigh, PriorityMedium, PriorityLow:
		return true
	}
	return false
}

// toCanonicalPriority converts a case-insensitive priority string to its canonical PriorityLevel.
// Returns an empty string if the input does not match any known priority.
func toCanonicalPriority(p PriorityLevel) PriorityLevel {
	if canonical, ok := priorityMap[strings.ToLower(string(p))]; ok {
		return canonical
	}
	return ""
}

// Complete marks a todo item as completed by its ID.
// It iterates through the list to find the matching todo and updates its `Completed` status.
// Returns an error if the todo with the given ID is not found.
func (tl *TodoList) Complete(id int) error {
	// Iterate over the slice of todos using index `i`.
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			// If the ID matches, mark the todo as completed.
			tl.Todos[i].Completed = true
			PrintUserMessage(fmt.Sprintf("🎉 Completed todo #%d: \"%s\"", tl.Todos[i].ID, tl.Todos[i].Task))
			return nil // Return nil on success.
		}
	}
	// If no todo with the given ID is found after iterating, return an error.
	return fmt.Errorf("todo with ID %d not found", id)
}

// Uncomplete marks a todo item as incomplete by its ID.
// It iterates through the list to find the matching todo and updates its `Completed` status to false.
// Returns an error if the todo with the given ID is not found.
func (tl *TodoList) Uncomplete(id int) error {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			// If the ID matches, mark the todo as incomplete.
			tl.Todos[i].Completed = false
			PrintUserMessage(fmt.Sprintf("🔄 Uncompleted todo #%d: \"%s\"", tl.Todos[i].ID, tl.Todos[i].Task))
			return nil // Return nil on success.
		}
	}
	// If no todo with the given ID is found after iterating, return an error.
	return fmt.Errorf("todo with ID %d not found", id)
}

// Delete removes a todo item from the TodoList by its ID.
// It iterates through the list, finds the matching todo, and removes it by creating a new slice
// that excludes the deleted item. Returns the deleted Todo and an error if not found.
func (tl *TodoList) Delete(id int) (Todo, error) {
	// Iterate over the slice of todos with both index `i` and `todo` value.
	for i, todo := range tl.Todos {
		if todo.ID == id {
			// If the ID matches, remove the todo from the slice.
			// This is done by appending the slice before the item to the slice after the item.
			tl.Todos = append(tl.Todos[:i], tl.Todos[i+1:]...)
			PrintUserMessage(fmt.Sprintf("🗑️ Deleted todo #%d: \"%s\"", todo.ID, todo.Task))
			return todo, nil // Return the deleted todo and nil on success.
		}
	}
	// If no todo with the given ID is found, return an error.
	return Todo{}, fmt.Errorf("todo with ID %d not found", id)
}

// ListOptions defines parameters for filtering and sorting todos.
type ListOptions struct {
	FilterStatus   string        // "all", "completed", "incomplete"
	FilterPriority PriorityLevel // Specific priority (e.g., "high")
	FilterTags     []string      // Tags to filter by
	SortBy         string        // "id", "task", "created_at", "due_date", "priority"
	SortOrder      string        // "asc" (ascending) or "desc" (descending)
}

// List prints all todo items in the TodoList to the console, applying optional filters and sorting.
func (tl *TodoList) List(options ListOptions) {
	filteredTodos := []Todo{}
	for _, todo := range tl.Todos {
		match := true

		// Filter by status
		if options.FilterStatus == "completed" && !todo.Completed {
			match = false
		}
		if options.FilterStatus == "incomplete" && todo.Completed {
			match = false
		}

		// Filter by priority
		// Normalize the filter priority for case-insensitive comparison
		canonicalFilterPriority := toCanonicalPriority(options.FilterPriority)
		if canonicalFilterPriority != "" && todo.Priority != canonicalFilterPriority {
			match = false
		}

		// Filter by tags
		if len(options.FilterTags) > 0 {
			hasTag := false
			for _, filterTag := range options.FilterTags {
				for _, todoTag := range todo.Tags {
					if strings.ToLower(filterTag) == strings.ToLower(todoTag) {
						hasTag = true
						break
					}
				}
				if hasTag {
					break
				}
			}
			if !hasTag {
				match = false
			}
		}

		if match {
			filteredTodos = append(filteredTodos, todo)
		}
	}

	// Sort todos if sortBy is specified.
	if options.SortBy != "" {
		sort.Slice(filteredTodos, func(i, j int) bool {
			a, b := filteredTodos[i], filteredTodos[j]
			var less bool

			switch options.SortBy {
			case "id":
				less = a.ID < b.ID
			case "task":
				less = strings.ToLower(a.Task) < strings.ToLower(b.Task)
			case "created_at":
				less = a.CreatedAt.Before(b.CreatedAt)
			case "due_date":
				// Handle nil DueDate for sorting
				if options.SortOrder == "desc" {
					// Descending order: later dates first, nil dates last
					if a.DueDate == nil && b.DueDate == nil {
						return false // Treat as equal
					} else if a.DueDate == nil {
						return false // Nil due dates come after non-nil, so 'a' is not 'less' than 'b'
					} else if b.DueDate == nil {
						return true // Non-nil due dates come before nil, so 'a' is 'less' than 'b'
					} else {
						return a.DueDate.After(*b.DueDate) // For desc, 'a' comes before 'b' if 'a' is after 'b'
					}
				} else { // Ascending order
					// Ascending order: earlier dates first, nil dates last
					if a.DueDate == nil && b.DueDate == nil {
						return false // Treat as equal
					} else if a.DueDate == nil {
						return false // Nil due dates come after non-nil, so 'a' is not 'less' than 'b'
					} else if b.DueDate == nil {
						return true // Non-nil due dates come before nil, so 'a' is 'less' than 'b'
					} else {
						return a.DueDate.Before(*b.DueDate) // For asc, 'a' comes before 'b' if 'a' is before 'b'
					}
				}
			case "priority":
				// Simple alphabetical sort for priority for now; can be enhanced with custom order.
				less = strings.ToLower(string(a.Priority)) < strings.ToLower(string(b.Priority))
			default:
				// Default sort by ID if sortBy is unknown
				less = a.ID < b.ID
			}

			if options.SortOrder == "desc" {
				return !less
			}
			return less
		})
	}

	// Print the filtered and sorted todos.
	if len(filteredTodos) == 0 {
		PrintUserMessage("✨ No todos found matching the criteria.")
		return
	}

	PrintUserMessage("📋 Your Todos:")
	for _, todo := range filteredTodos {
		status := "[ ]"
		if todo.Completed {
			status = "[x]"
		}
		priorityStr := ""
		if todo.Priority != "" {
			// Capitalize the first letter for display
			priorityStr = fmt.Sprintf(" (Priority: %s)", strings.Title(string(todo.Priority)))
		}
		dueDateStr := ""
		if todo.DueDate != nil {
			dueDateStr = fmt.Sprintf(" (Due: %s)", todo.DueDate.Format("2006-01-02"))
		}
		tagsStr := ""
		if len(todo.Tags) > 0 {
			tagsStr = fmt.Sprintf(" [Tags: %s]", strings.Join(todo.Tags, ", "))
		}
		// Use PrintUserMessage for consistent output, including emojis
		PrintUserMessage(fmt.Sprintf("%s %d. %s%s%s%s (Created: %s)", status, todo.ID, todo.Task, priorityStr, dueDateStr, tagsStr, todo.CreatedAt.Format("2006-01-02 15:04")))
	}
}

// SearchTasks finds todo items whose task description or tags contain the given query string.
// The search is case-insensitive.
func (tl *TodoList) SearchTasks(query string) *TodoList {
	matchedTodos := NewTodoList()
	lowerQuery := strings.ToLower(query)

	for _, todo := range tl.Todos {
		// Check if the task description contains the query.
		if strings.Contains(strings.ToLower(todo.Task), lowerQuery) {
			matchedTodos.Todos = append(matchedTodos.Todos, todo)
			continue // Move to the next todo once a match is found in task description.
		}
		// Check if any tag contains the query.
		for _, tag := range todo.Tags {
			if strings.Contains(strings.ToLower(tag), lowerQuery) {
				matchedTodos.Todos = append(matchedTodos.Todos, todo)
				break // Move to the next todo once a match is found in tags.
			}
		}
	}

	return matchedTodos
}

// ClearCompleted removes all completed todo items from the list.
func (tl *TodoList) ClearCompleted() {
	var activeTodos []Todo
	for _, todo := range tl.Todos {
		if !todo.Completed {
			activeTodos = append(activeTodos, todo)
		}
	}
	// Check if any todos were actually removed.
	if len(tl.Todos) > len(activeTodos) {
		PrintUserMessage(fmt.Sprintf("🧹 Cleared %d completed todos.", len(tl.Todos)-len(activeTodos)))
		tl.Todos = activeTodos
	} else {
		PrintUserMessage("No completed todos to clear.")
	}
}

// EditTask updates the task description of an existing todo item.
// It takes the ID of the todo to edit and the new task description.
// Returns an error if the todo with the given ID is not found.
func (tl *TodoList) EditTask(id int, newTask string) error {
	for i := range tl.Todos {
		if tl.Todos[i].ID == id {
			// Update the task description.
			taskBefore := tl.Todos[i].Task
			tl.Todos[i].Task = newTask
			PrintUserMessage(fmt.Sprintf("✏️ Edited todo #%d. Old task: \"%s\", New task: \"%s\"", id, taskBefore, newTask))
			return nil // Return nil on success.
		}
	}
	// If no todo with the given ID is found, return an error.
	return fmt.Errorf("todo with ID %d not found", id)
}

// SaveToFile saves the current state of the TodoList to a JSON file.
// It marshals the `TodoList` struct into a pretty-printed JSON format and writes it to the specified file.
// Returns an error if marshaling or file writing fails.
func (tl *TodoList) SaveToFile(filename string) error {
	// Marshal the TodoList struct into JSON format with indentation.
	data, err := json.MarshalIndent(tl, "", "  ")
	LogError(err, "Failed to marshal todo list to JSON") // Uncommented LogError
	if err != nil {
		return fmt.Errorf("failed to save todos: %w", err)
	}

	// Write the JSON data to the specified file with read/write permissions for the owner.
	err = os.WriteFile(filename, data, 0644)
	LogError(err, fmt.Sprintf("Failed to write todo list to file %s", filename)) // Uncommented LogError
	if err != nil {
		return fmt.Errorf("failed to save todos to file: %w", err)
	}

	LogInfo(fmt.Sprintf("Todos saved to %s", filename)) // Uncommented LogInfo
	return nil                                          // Return nil on successful save.
}

// LoadFromFile loads a TodoList from a JSON file.
// It reads the file, unmarshals the JSON data into a `TodoList` struct.
// If the file does not exist, it returns a new empty `TodoList`.
// Returns an error if file reading or JSON unmarshaling fails.
func LoadFromFile(filename string) (*TodoList, error) {
	// Read the content of the specified file.
	data, err := os.ReadFile(filename)
	LogError(err, fmt.Sprintf("Failed to read todo list from file %s", filename)) // Uncommented LogError
	if err != nil {
		// If the file does not exist, create and return a new empty TodoList without an error.
		if os.IsNotExist(err) {
			LogInfo(fmt.Sprintf("⚠️ Todo file %s does not exist, creating new list.", filename)) // Uncommented LogInfo
			return NewTodoList(), nil
		}
		// For other file reading errors, return an error.
		return nil, fmt.Errorf("failed to load todos: %w", err)
	}

	// Create a new TodoList to unmarshal the data into.
	todoList := NewTodoList()
	// Unmarshal the JSON data from the file into the todoList struct.
	err = json.Unmarshal(data, todoList)
	LogError(err, "Failed to unmarshal todo list from JSON") // Uncommented LogError
	if err != nil {
		return nil, fmt.Errorf("failed to parse todos: %w", err)
	}

	LogInfo(fmt.Sprintf("Todos loaded from %s", filename)) // Uncommented LogInfo
	return todoList, nil                                   // Return the loaded todo list and nil on success.
}
