package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mohammadrabetian/ports/domain"
)

type PortHandlerFunc func(port *domain.Port) error

/* NOTES: optimizations (switch, label, cache) */
func ProcessJSONFile(filename string, portHandler PortHandlerFunc) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read the opening brace of the JSON object
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("failed to read the opening brace of JSON object: %w", err)
	}

	for {
		// Check if the next token is a closing brace, indicating the end of the JSON object
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read JSON token: %w", err)
		}
		if delim, ok := t.(json.Delim); ok && delim == '}' {
			break
		}

		// If not a closing brace, process the key and corresponding port data
		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("expected a string key, got %T", t)
		}

		var port domain.Port

		if err := decoder.Decode(&port); err != nil {
			return fmt.Errorf("failed to decode port data: %w", err)
		}

		/* NOTES: validations */
		// Set the port's ID using the JSON key
		port.ID = key

		// Call the provided port handler function
		if err := portHandler(&port); err != nil {
			return fmt.Errorf("failed to handle port: %w", err)
		}
	}

	return nil
}
