package util_test

import (
	"errors"
	"os"
	"testing"

	"github.com/mohammadrabetian/ports/domain"
	"github.com/mohammadrabetian/ports/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessJSONFile_Success(t *testing.T) {
	// Create a temporary JSON file
	tmpFile, err := os.CreateTemp("", "ports-test-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write JSON content to the temporary file
	fileContent := `{"ABC": {"name": "Test Port", "city": "Test City", "country": "Test Country"}}`
	_, err = tmpFile.Write([]byte(fileContent))
	require.NoError(t, err)
	tmpFile.Close()

	portHandlerCalled := false
	portHandler := func(port *domain.Port) error {
		assert.Equal(t, "ABC", port.ID)
		assert.Equal(t, "Test Port", port.Name)
		assert.Equal(t, "Test City", port.City)
		assert.Equal(t, "Test Country", port.Country)
		portHandlerCalled = true
		return nil
	}

	err = util.ProcessJSONFile(tmpFile.Name(), portHandler)
	assert.NoError(t, err)
	assert.True(t, portHandlerCalled, "portHandler function should be called")
}

func TestProcessJSONFile_InvalidJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ports-test-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write invalid JSON content to the temporary file
	fileContent := `{"ABC": {"name": "Test Port", "city": "Test City", "country": "Test Country"`
	_, err = tmpFile.Write([]byte(fileContent))
	require.NoError(t, err)
	tmpFile.Close()

	portHandler := func(port *domain.Port) error {
		return nil
	}

	err = util.ProcessJSONFile(tmpFile.Name(), portHandler)
	assert.Error(t, err)
}

func TestProcessJSONFile_EmptyJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ports-test-empty-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write empty JSON content to the temporary file
	fileContent := `{}`
	_, err = tmpFile.Write([]byte(fileContent))
	require.NoError(t, err)
	tmpFile.Close()

	portHandlerCalled := false
	portHandler := func(port *domain.Port) error {
		portHandlerCalled = true
		return nil
	}

	err = util.ProcessJSONFile(tmpFile.Name(), portHandler)
	assert.NoError(t, err)
	assert.False(t, portHandlerCalled, "portHandler function should not be called")
}

func TestProcessJSONFile_PortHandlerError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "ports-test-handler-error-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	fileContent := `{"ABC": {"name": "Test Port", "city": "Test City", "country": "Test Country"}}`
	_, err = tmpFile.Write([]byte(fileContent))
	require.NoError(t, err)
	tmpFile.Close()

	// Prepare the port handler function that always returns an error
	portHandler := func(port *domain.Port) error {
		return errors.New("port handler error")
	}

	err = util.ProcessJSONFile(tmpFile.Name(), portHandler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to handle port")
	assert.Contains(t, err.Error(), "port handler error")
}
