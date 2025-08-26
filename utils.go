package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log" // Package for logging functionality
	"os"  // Package for operating system functionalities, used here for stderr
	"time"
)

// Duration is a custom type that allows time.Duration to be marshaled/unmarshaled
// from JSON as a human-readable string (e.g., "30s", "1m").
type Duration time.Duration

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (d *Duration) UnmarshalText(text []byte) error {
	parsed, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration(parsed)
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// SetupLogger configures the application-wide logger.
// It sets the output destination to standard error (os.Stderr) and defines
// the logging flags to include date, time, and source file information.
// Optionally, logs can also be directed to a specified file.
func SetupLogger(logFilePath string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if logFilePath != "" {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("ERROR: Failed to open log file %s: %v\n", logFilePath, err)
			log.SetOutput(os.Stderr) // Fallback to stderr if file logging fails
			return
		}
		// Create a multi-writer to write to both stderr and the file.
		mw := io.MultiWriter(os.Stderr, file)
		log.SetOutput(mw)
	} else {
		log.SetOutput(os.Stderr) // Default to stderr
	}
}

// LogError logs an error message with a specified context.
// It takes an error object and a descriptive message. If the error is not nil,
// it prints a formatted error message including the custom message and the error details.
func LogError(err error, message string) {
	if err != nil {
		// Use Printf to format the error message, including the custom message and the error itself.
		log.Printf("ERROR: %s: %v\n", message, err)
	}
}

// LogInfo logs an informational message.
// It takes a descriptive message and prints it as an informational log entry.
func LogInfo(message string) {
	// Use Printf to format the informational message.
	log.Printf("INFO: %s\n", message)
}

// LogWarning logs a warning message.
// It takes a descriptive message and prints it as a warning log entry.
func LogWarning(message string) {
	log.Printf("WARNING: %s\n", message)
}

// PrintUserMessage prints messages directly to standard output, without any logger prefixes.
// This is intended for direct user feedback in the CLI.
func PrintUserMessage(message string) {
	fmt.Println(message)
}

// Config holds the application's configurable settings.
type Config struct {
	DataFile         string   `json:"data_file"`
	AutoSaveInterval Duration `json:"auto_save_interval"` // Use custom Duration type
	LogFilePath      string   `json:"log_file_path"`
}

// DefaultConfig returns a new Config with default values.
func DefaultConfig() Config {
	return Config{
		DataFile:         "todos.json",
		AutoSaveInterval: Duration(1 * time.Minute), // Cast to custom Duration type
		LogFilePath:      "",                        // Default to no log file (stdout/stderr only)
	}
}

// LoadConfig loads configuration from a JSON file. If the file does not exist,
// it creates a default configuration file.
func LoadConfig(configPath string) (Config, error) {
	config := DefaultConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If config file doesn't exist, create a default one.
			LogInfo(fmt.Sprintf("Config file %s not found, creating default.", configPath))
			err = SaveConfig(config, configPath)
			if err != nil {
				LogError(err, fmt.Sprintf("Failed to save default config to %s", configPath))
				return config, fmt.Errorf("failed to create default config: %w", err)
			}
			return config, nil
		}
		LogError(err, fmt.Sprintf("Failed to read config file %s", configPath))
		return config, fmt.Errorf("failed to load config: %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		LogError(err, fmt.Sprintf("Failed to parse config file %s", configPath))
		return config, fmt.Errorf("failed to parse config: %w", err)
	}

	LogInfo(fmt.Sprintf("Config loaded from %s.", configPath))
	return config, nil
}

// SaveConfig saves the given Config to a JSON file.
func SaveConfig(config Config, configPath string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		LogError(err, "Failed to marshal config to JSON")
		return fmt.Errorf("failed to save config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		LogError(err, fmt.Sprintf("Failed to write config to file %s", configPath))
		return fmt.Errorf("failed to save config to file: %w", err)
	}

	LogInfo(fmt.Sprintf("Config saved to %s.", configPath))
	return nil
}
