package main

import (
	"bytes"   // New import for bytes.Buffer
	"io"      // Package for input/output operations, used for capturing stdout
	"log"     // Package for logging, used for capturing log output
	"os"      // Package for operating system functionalities, used for file removal
	"reflect" // Package for reflection, used for deep comparison of structs
	"strings" // Package for string manipulation, used for capturing and checking output
	"testing" // Package for writing automated tests
	"time"    // Package for time-related operations, used for `time.Duration` and `time.Sleep`
)

// TestNewTodoList verifies that NewTodoList initializes an empty list with the correct NextID.
func TestNewTodoList(t *testing.T) {
	// Call the function under test.
	tl := NewTodoList()

	// Assert that the returned TodoList is not nil.
	if tl == nil {
		t.Error("NewTodoList() returned nil")
	}

	// Assert that the Todos slice is empty.
	if len(tl.Todos) != 0 {
		t.Errorf("NewTodoList() should return an empty list, got %d todos", len(tl.Todos))
	}

	// Assert that NextID is initialized to 1.
	if tl.NextID != 1 {
		t.Errorf("NewTodoList() should initialize NextID to 1, got %d", tl.NextID)
	}
}

func TestAdd(t *testing.T) {
	tl := NewTodoList()
	task := "Buy groceries"
	priority := PriorityLevel("high")
	dueDateStr := "2024-12-25"
	parsedDate, _ := time.Parse("2006-01-02", dueDateStr)
	dueDate := &parsedDate
	tags := []string{"personal", "shopping"}

	tl.Add(task, priority, dueDate, tags)

	if len(tl.Todos) != 1 {
		t.Errorf("Add() failed, expected 1 todo, got %d", len(tl.Todos))
	}
	if tl.Todos[0].Task != task {
		t.Errorf("Add() failed, expected task %s, got %s", task, tl.Todos[0].Task)
	}
	if tl.Todos[0].ID != 1 {
		t.Errorf("Add() failed, expected ID 1, got %d", tl.Todos[0].ID)
	}
	if tl.Todos[0].Priority != priority {
		t.Errorf("Add() failed, expected priority %s, got %s", priority, tl.Todos[0].Priority)
	}

	// Test case-insensitive priority
	caseInsensitivePriority := PriorityLevel("hIgH")
	tl.Add("Case-insensitive task", caseInsensitivePriority, nil, nil)
	if len(tl.Todos) != 2 {
		t.Errorf("Add() with case-insensitive priority failed, expected 2 todos, got %d", len(tl.Todos))
	}
	if tl.Todos[1].Priority != PriorityLevel("high") {
		t.Errorf("Add() with case-insensitive priority failed, expected priority %s, got %s", PriorityLevel("high"), tl.Todos[1].Priority)
	}

	if !tl.Todos[0].DueDate.Equal(*dueDate) {
		t.Errorf("Add() failed, expected due date %v, got %v", *dueDate, *tl.Todos[0].DueDate)
	}
	if !reflect.DeepEqual(tl.Todos[0].Tags, tags) {
		t.Errorf("Add() failed, expected tags %v, got %v", tags, tl.Todos[0].Tags)
	}
	if tl.NextID != 3 {
		t.Errorf("Add() failed, expected NextID 3, got %d", tl.NextID)
	}
}

func TestAdd_InvalidPriority(t *testing.T) {
	tl := NewTodoList()
	task := "Task with invalid priority"
	invalidPriority := PriorityLevel("Urgent") // An invalid priority string
	dueDateStr := "2024-12-25"
	parsedDate, _ := time.Parse("2006-01-02", dueDateStr)
	dueDate := &parsedDate
	tags := []string{"test"}

	// Capture log output to check for warning message
	oldOutput := log.Writer()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(oldOutput) // Restore original output after test

	tl.Add(task, invalidPriority, dueDate, tags)

	if len(tl.Todos) != 1 {
		t.Errorf("Add() failed, expected 1 todo, got %d", len(tl.Todos))
	}
	if tl.Todos[0].Priority != PriorityLevel("medium") {
		t.Errorf("Add() with invalid priority failed, expected priority %s, got %s", PriorityLevel("medium"), tl.Todos[0].Priority)
	}

	expectedWarning := "WARNING: Invalid priority level 'Urgent' for task 'Task with invalid priority'. Defaulting to Medium."
	if !strings.Contains(buf.String(), expectedWarning) {
		t.Errorf("Expected warning not found in log output. Got: %s", buf.String())
	}

	if tl.NextID != 2 {
		t.Errorf("Add() failed, expected NextID 2, got %d", tl.NextID)
	}
}

func TestComplete(t *testing.T) {
	tl := NewTodoList()
	tl.Add("Task 1", PriorityLevel("medium"), nil, nil)
	tl.Add("Task 2", PriorityLevel("medium"), nil, nil)

	err := tl.Complete(1)
	if err != nil {
		t.Errorf("Complete() failed: %v", err)
	}
	if !tl.Todos[0].Completed {
		t.Error("Complete() failed, Task 1 should be completed")
	}

	err = tl.Complete(99)
	if err == nil {
		t.Error("Complete() should return an error for non-existent ID")
	}
}

func TestUncomplete(t *testing.T) {
	tl := NewTodoList()
	tl.Add("Task 1", PriorityLevel("medium"), nil, nil)
	tl.Complete(1) // Mark as complete first

	err := tl.Uncomplete(1)
	if err != nil {
		t.Errorf("Uncomplete() failed: %v", err)
	}
	if tl.Todos[0].Completed {
		t.Error("Uncomplete() failed, Task 1 should be incomplete")
	}

	err = tl.Uncomplete(99)
	if err == nil {
		t.Error("Uncomplete() should return an error for non-existent ID")
	}
}

func TestDelete(t *testing.T) {
	tl := NewTodoList()
	tl.Add("Task 1", PriorityLevel("medium"), nil, nil)
	tl.Add("Task 2", PriorityLevel("medium"), nil, nil)
	tl.Add("Task 3", PriorityLevel("medium"), nil, nil)

	deletedTodo, err := tl.Delete(2)
	if err != nil {
		t.Errorf("Delete() failed: %v", err)
	}

	if deletedTodo.ID != 2 || deletedTodo.Task != "Task 2" {
		t.Errorf("Delete() returned incorrect todo. Expected ID 2, Task \"Task 2\", got ID %d, Task \"%s\"", deletedTodo.ID, deletedTodo.Task)
	}
	if len(tl.Todos) != 2 {
		t.Errorf("Delete() failed, expected 2 todos, got %d", len(tl.Todos))
	}

	if tl.Todos[0].ID != 1 || tl.Todos[1].ID != 3 {
		t.Errorf("Delete() failed, incorrect todos remaining")
	}

	_, err = tl.Delete(99)
	if err == nil {
		t.Error("Delete() should return an error for non-existent ID")
	}
}

func TestEditTask(t *testing.T) {
	tl := NewTodoList()
	tl.Add("Original Task", PriorityLevel("medium"), nil, nil)

	newTask := "Edited Task Description"
	err := tl.EditTask(1, newTask)
	if err != nil {
		t.Errorf("EditTask() failed: %v", err)
	}
	if tl.Todos[0].Task != newTask {
		t.Errorf("EditTask() failed, expected task %s, got %s", newTask, tl.Todos[0].Task)
	}

	err = tl.EditTask(99, "Nonexistent Task")
	if err == nil {
		t.Error("EditTask() should return an error for non-existent ID")
	}
}

func TestClearCompleted(t *testing.T) {
	tl := NewTodoList()
	tl.Add("Task 1", PriorityLevel("medium"), nil, nil)
	tl.Add("Task 2", PriorityLevel("medium"), nil, nil)
	tl.Add("Task 3", PriorityLevel("medium"), nil, nil)

	tl.Complete(1)
	tl.Complete(3)

	// Clear completed todos
	tl.ClearCompleted()

	if len(tl.Todos) != 1 {
		t.Errorf("ClearCompleted() failed, expected 1 todo, got %d", len(tl.Todos))
	}
	if tl.Todos[0].ID != 2 {
		t.Errorf("ClearCompleted() failed, expected todo #2, got #%d", tl.Todos[0].ID)
	}

	// Test clearing when no todos are completed
	tl2 := NewTodoList()
	tl2.Add("Task A", PriorityLevel("medium"), nil, nil)
	tl2.ClearCompleted()
	if len(tl2.Todos) != 1 {
		t.Errorf("ClearCompleted() failed when no todos are complete, expected 1, got %d", len(tl2.Todos))
	}

	// Test clearing an empty list
	tl3 := NewTodoList()
	tl3.ClearCompleted()
	if len(tl3.Todos) != 0 {
		t.Errorf("ClearCompleted() failed for empty list, expected 0, got %d", len(tl3.Todos))
	}
}

func TestSearchTasks(t *testing.T) {
	tl := NewTodoList()
	parsedDate1, _ := time.Parse("2006-01-02", "2024-01-01")
	parsedDate2, _ := time.Parse("2006-01-02", "2024-01-02")

	tl.Add("Buy groceries", PriorityLevel("high"), &parsedDate1, []string{"personal", "shopping"})
	tl.Add("Work on report", PriorityLevel("medium"), &parsedDate2, []string{"work", "urgent"})
	tl.Add("Call John", PriorityLevel("low"), nil, []string{"personal"})
	tl.Add("Review code", PriorityLevel("high"), nil, []string{"work", "code"})

	// Test search by task description
	results1 := tl.SearchTasks("report")
	if len(results1.Todos) != 1 || results1.Todos[0].Task != "Work on report" {
		t.Errorf("SearchTasks(\"report\") failed. Expected 1 task 'Work on report', got %d tasks", len(results1.Todos))
	}

	// Test search by tag
	results2 := tl.SearchTasks("work")
	if len(results2.Todos) != 2 {
		t.Errorf("SearchTasks(\"work\") failed. Expected 2 tasks, got %d", len(results2.Todos))
	}

	// Test case-insensitive search
	results3 := tl.SearchTasks("groceries")
	if len(results3.Todos) != 1 || results3.Todos[0].Task != "Buy groceries" {
		t.Errorf("SearchTasks(\"groceries\") failed. Expected 1 task 'Buy groceries', got %d tasks", len(results3.Todos))
	}

	// Test no results
	results4 := tl.SearchTasks("nonexistent")
	if len(results4.Todos) != 0 {
		t.Errorf("SearchTasks(\"nonexistent\") failed. Expected 0 tasks, got %d", len(results4.Todos))
	}

	// Test empty query (should return all)
	results5 := tl.SearchTasks("")
	if len(results5.Todos) != 4 {
		t.Errorf("SearchTasks(\"\") failed. Expected 4 tasks, got %d", len(results5.Todos))
	}
}

func TestListWithFilteringAndSorting(t *testing.T) {
	tl := NewTodoList()
	parsedDate1, _ := time.Parse("2006-01-02", "2024-01-01")
	parsedDate2, _ := time.Parse("2006-01-02", "2024-01-02")

	tl.Add("Task B Low", PriorityLevel("low"), &parsedDate2, []string{"personal"})
	tl.Add("Task A High", PriorityLevel("high"), &parsedDate1, []string{"work"})
	tl.Add("Task C Med", PriorityLevel("medium"), nil, []string{"personal", "urgent"})

	// Test filter by status (incomplete)
	options1 := ListOptions{FilterStatus: "incomplete"}
	out := captureOutput(func() { tl.List(options1) })
	if !strings.Contains(out, "Task B Low") || !strings.Contains(out, "Task A High") || !strings.Contains(out, "Task C Med") {
		t.Errorf("List with FilterStatus incomplete failed: %s", out)
	}

	// Mark one as complete
	tl.Complete(1)
	options2 := ListOptions{FilterStatus: "completed"}
	out = captureOutput(func() { tl.List(options2) })
	if !strings.Contains(out, "Task B Low") || strings.Contains(out, "Task A High") || strings.Contains(out, "Task C Med") {
		t.Errorf("List with FilterStatus completed failed: %s", out)
	}

	// Test filter by priority
	options3 := ListOptions{FilterPriority: PriorityLevel("hIgH")}
	out = captureOutput(func() { tl.List(options3) })
	if !strings.Contains(out, "Task A High") || strings.Contains(out, "Task B Low") {
		t.Errorf("List with FilterPriority high failed: %s", out)
	}

	// Test filter by tags
	options4 := ListOptions{FilterTags: []string{"urgent"}}
	out = captureOutput(func() { tl.List(options4) })
	if !strings.Contains(out, "Task C Med") || strings.Contains(out, "Task A High") {
		t.Errorf("List with FilterTags urgent failed: %s", out)
	}

	// Test sort by created_at asc
	options5 := ListOptions{SortBy: "created_at", SortOrder: "asc"}
	out = captureOutput(func() { tl.List(options5) })
	expectedOrder5 := []string{"Task B Low", "Task A High", "Task C Med"}
	if !checkOrder(out, expectedOrder5) {
		t.Errorf("List with SortBy created_at asc failed: %s", out)
	}

	// Test sort by due_date desc
	options6 := ListOptions{SortBy: "due_date", SortOrder: "desc"}
	out = captureOutput(func() { tl.List(options6) })
	expectedOrder6 := []string{"Task B Low", "Task A High", "Task C Med"} // CMed has nil due date, comes last
	if !checkOrder(out, expectedOrder6) {
		t.Errorf("List with SortBy due_date desc failed: %s", out)
	}

	// Test sort by priority asc (alphabetical)
	options7 := ListOptions{SortBy: "priority", SortOrder: "asc"}
	out = captureOutput(func() { tl.List(options7) })
	expectedOrder7 := []string{"Task A High", "Task B Low", "Task C Med"} // high, low, medium alphabetically
	if !checkOrder(out, expectedOrder7) {
		t.Errorf("List with SortBy priority asc failed: %s", out)
	}

}

// Helper to capture fmt.Println output for testing.
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout
	return string(out)
}

// Helper to check the order of tasks in the output string.
func checkOrder(output string, expectedOrder []string) bool {
	outputLines := strings.Split(strings.TrimSpace(output), "\n")
	// Skip the header "📋 Your Todos:" if present
	if len(outputLines) > 0 && strings.HasPrefix(outputLines[0], "📋 Your Todos:") {
		outputLines = outputLines[1:]
	}

	if len(outputLines) != len(expectedOrder) {
		return false
	}

	for i, expectedTask := range expectedOrder {
		// Extract just the task description from the output line for comparison.
		// Example line format: "[ ] 1. Task B Low (Priority: low) (Due: 2024-01-02) [Tags: personal] (Created: 2025-08-24 22:20)"
		// We want to extract "Task B Low".
		line := outputLines[i]
		// Find the start of the task description (after ID and ". ")
		idx := strings.Index(line, ". ")
		if idx == -1 || idx+2 >= len(line) {
			return false // Malformed line, cannot extract task
		}
		line = line[idx+2:] // Remove status and ID part

		// Find the end of the task description (before first parenthesis or bracket, or end of line)
		endIdx := len(line)
		if pIdx := strings.Index(line, " ("); pIdx != -1 && pIdx < endIdx {
			endIdx = pIdx
		}
		if bIdx := strings.Index(line, " ["); bIdx != -1 && bIdx < endIdx {
			endIdx = bIdx
		}
		taskInOutput := strings.TrimSpace(line[:endIdx])

		if taskInOutput != expectedTask {
			return false
		}
	}
	return true
}

func TestSaveAndLoad(t *testing.T) {
	// Create a temporary file for testing
	testFilename := "test_todos.json"
	defer os.Remove(testFilename) // Clean up the file after the test

	// Create a new todo list and add some items
	tl1 := NewTodoList()
	parsedDate, _ := time.Parse("2006-01-02", "2024-03-15")
	tl1.Add("Task A with details", PriorityLevel("high"), &parsedDate, []string{"work", "project"})
	tl1.Add("Task B simple", PriorityLevel("low"), nil, nil)

	// Save the list to file
	err := tl1.SaveToFile(testFilename)
	if err != nil {
		t.Fatalf("SaveToFile() failed: %v", err)
	}

	// Load the list from file into a new TodoList
	tl2, err := LoadFromFile(testFilename)
	if err != nil {
		t.Fatalf("LoadFromFile() failed: %v", err)
	}

	// Verify the loaded list matches the original
	if len(tl1.Todos) != len(tl2.Todos) {
		t.Errorf("Loaded list has different number of todos, expected %d, got %d", len(tl1.Todos), len(tl2.Todos))
	}

	for i := range tl1.Todos {
		// Deep compare requires reflecting on all fields, including pointers for DueDate
		if tl1.Todos[i].ID != tl2.Todos[i].ID ||
			tl1.Todos[i].Task != tl2.Todos[i].Task ||
			tl1.Todos[i].Completed != tl2.Todos[i].Completed ||
			tl1.Todos[i].Priority != tl2.Todos[i].Priority ||
			!reflect.DeepEqual(tl1.Todos[i].DueDate, tl2.Todos[i].DueDate) ||
			!reflect.DeepEqual(tl1.Todos[i].Tags, tl2.Todos[i].Tags) {
			t.Errorf("Loaded todo item at index %d does not match original. Expected: %+v, Got: %+v", i, tl1.Todos[i], tl2.Todos[i])
		}
	}

	// Test loading from a non-existent file (should return an empty list)
	os.Remove(testFilename) // Ensure the file doesn't exist
	tl3, err := LoadFromFile(testFilename)
	if err != nil {
		t.Fatalf("LoadFromFile() failed for non-existent file: %v", err)
	}
	if len(tl3.Todos) != 0 {
		t.Errorf("LoadFromFile() for non-existent file should return empty list, got %d todos", len(tl3.Todos))
	}
}

func TestAutoSave(t *testing.T) {
	testFilename := "test_autosave.json"
	defer os.Remove(testFilename)

	tl := NewTodoList()

	// Start auto-save with a short interval for testing
	interval := 100 * time.Millisecond
	StartAutoSave(tl, testFilename, interval)

	// Add a task and wait for a bit longer than the interval
	tl.Add("Auto-save task with priority", PriorityLevel("medium"), nil, []string{"auto"})
	time.Sleep(interval + (50 * time.Millisecond))

	// Load the file to check if the task was saved
	loadedTl, err := LoadFromFile(testFilename)
	if err != nil {
		t.Fatalf("Failed to load file after auto-save: %v", err)
	}
	if len(loadedTl.Todos) != 1 || loadedTl.Todos[0].Task != "Auto-save task with priority" || loadedTl.Todos[0].Priority != PriorityLevel("medium") || !reflect.DeepEqual(loadedTl.Todos[0].Tags, []string{"auto"}) {
		t.Errorf("Auto-save failed, expected 'Auto-save task with priority' with details to be saved")
	}
}
