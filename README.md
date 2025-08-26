# CLI Todo List Application (Golang)

This is a command-line interface (CLI) application for managing a todo list, built with Golang. It emphasizes modularity, clean code principles, robust error handling, logging, and data persistence.

## Features

*   **Enhanced Todo Model:** Tasks now support `Priority`, `Due Date`, and `Tags`.
*   **Modular Design:** Code is organized into separate files based on their responsibilities, enhancing readability and maintainability.
*   **Error Handling & Logging:** Comprehensive error handling and application-wide logging provide better diagnostics and robustness. Logs can optionally be directed to a file.
*   **JSON Persistence:** Todo list data is automatically saved to and loaded from a `todos.json` file.
*   **Auto-Save Goroutine:** A background goroutine periodically saves the todo list, preventing data loss.
*   **Interactive Mode:** A continuous interactive mode allows users to manage todos without restarting the application for each command.
*   **New Commands:**`clear-completed`, `search`, `uncomplete`, `undo`.
*   **Advanced Listing:** The `list` command in **single-command mode** supports filtering by status, priority, and tags, as well as sorting by various fields.
*   **Confirmation Prompts:** Destructive actions like `delete` and `clear-completed` require user confirmation.
*   **Configurable Settings:** Application settings are managed via a `config.json` file.
*   **Unit Tests:** Core functionalities are covered by unit tests to ensure correctness.

## Project Structure

-   `cli/todo/main.go`: The application's entry point. Initializes the logger, loads/saves the todo list, starts the auto-save goroutine, and delegates command handling.
-   `cli/todo/models.go`: Defines the `Todo` and `TodoList` data structures and their core methods (add, complete, delete, list with options, save/load, edit, clear completed, search, uncomplete).
-   `cli/todo/utils.go`: Provides utility functions for configuring and using the application's logger, and handles configuration file loading/saving.
-   `cli/todo/autosave.go`: Implements the `StartAutoSave` function, which runs a goroutine for periodic data persistence.
-   `cli/todo/cli.go`: Handles all command-line argument parsing and manages the interactive user interface, including parsing complex interactive commands and providing confirmation prompts.
-   `cli/todo/models_test.go`: Contains comprehensive unit tests for all `TodoList` functionalities, including new features.
-   `cli/todo/go.mod`: Go module definition file for dependency management.
-   `cli/todo/todos.json`: (Created dynamically) Stores your todo list data in JSON format.
-   `cli/todo/config.json`: (Created dynamically) Stores application configuration settings.

## Getting Started

### Prerequisites

Ensure you have Go installed on your system. You can download it from [golang.org](https://golang.org/dl/).

### Running the Application

1.  **Navigate to the project directory:**
    Open your terminal or command prompt and change your current directory to the `cli/todo` folder:
    ```bash
    cd cli/todo
    ```

2.  **Initialize Go Module (if not already done):**
    If you haven't run `go mod init` before, you need to initialize the Go module. This creates the `go.mod` file for dependency management:
    ```bash
    go mod init todo
    ```
    *(You only need to run this command once per project.)*

3.  **Run the application:**
    To run the application, use `go run .` followed by your desired command-line flags. The `.` tells Go to compile and run all `.go` files in the current directory that belong to the `main` package.

    #### Single-Command Mode

    Execute a single action and then the application exits:

    *   **Add a new todo:**
        ```bash
        go run . -add "Learn Go modules" # Note: Priority, Due Date, Tags not supported via single -add flag currently.
        ```
    *   **Mark a todo as complete:**
        ```bash
        go run . -complete 1
        ```
    *   **Mark a todo as incomplete:**
        ```bash
        go run . -uncomplete 1
        ```
    *   **Delete a todo:** (Requires confirmation)
        ```bash
        go run . -delete 1
        ```
    *   **Clear all completed todos:** (Requires confirmation)
        ```bash
        go run . -clear-completed
        ```
    *   **List all todos (with filtering and sorting options):**
        ```bash
        go run . -list -filter-status incomplete -filter-priority high -filter-tags work,urgent -sort-by due_date -sort-order desc
        go run . -list # Simple list
        ```
    *   **View all available options/flags:**
        ```bash
        go run .
        ```

    #### Interactive Mode

    Run the application in a continuous interactive session. This mode is best for seeing the auto-save and logging features in action over time, and supports all enhanced commands.

    ```bash
    go run . -interactive
    ```

    Once in interactive mode, you will see a `>` prompt. Type your commands:

    *   `add Finish README -p high -d 2024-04-30 -t docs,urgent`
    *   `edit 1 "Refined README content"`
    *   `search "README"`
    *   `complete 1`
    *   `uncomplete 1`
    *   `delete 2` (Requires confirmation)
    *   `clear-completed` (Requires confirmation)
    *   `undo` (Undoes the last `add`, `complete`, `delete`, or `uncomplete`)
    *   `help` (for a list of interactive commands)
    *   `exit` (to quit interactive mode)

    *Note: In interactive mode, auto-save will periodically save your list in the background. Advanced listing options (filtering and sorting) are only available in single-command mode.* 

## Configuration

The application uses a `config.json` file for settings. If this file does not exist, a default one will be created when the application starts.

Example `config.json`:

```json
{
  "data_file": "todos.json",
  "auto_save_interval": "1m0s",
  "log_file_path": "app.log"
}
```

-   `data_file`: The name of the JSON file where todos are stored.
-   `auto_save_interval`: The interval at which the todo list is automatically saved (e.g., "1m0s" for 1 minute).
-   `log_file_path`: Optional. If set, application logs will be written to this file in addition to `stderr`.

## Running Tests

To run the unit tests for the application:

1.  **Navigate to the project directory:**
    ```bash
    cd cli/todo
    ```
2.  **Execute the tests:**
    ```bash
    go test
    ```
