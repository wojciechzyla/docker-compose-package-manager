package src

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestExecuteTemplate(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse("{{.Name}}"))

	tests := []struct {
		name       string
		values     map[string]interface{}
		expectFile bool
		expectErr  bool
	}{
		{"ValidTemplate", map[string]interface{}{"Name": "test_name"}, true, false},
		{"EmptyTemplate", map[string]interface{}{"Name": ""}, false, false},
		{"NilTemplate", nil, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_output")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			err = executeTemplate(tmpFile.Name(), tt.values, tmpl)
			if (err != nil) != tt.expectErr {
				t.Errorf("executeTemplate() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			content, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			if tt.expectFile {
				if len(content) == 0 {
					t.Errorf("Expected file to have content, but it was empty")
				}
			} else {
				if len(content) > 0 {
					t.Errorf("Expected file to be empty, but it had content: %s", content)
				}
			}
		})
	}
}

func TestValuesFromYamlFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expectErr   bool
		expected    map[string]interface{}
	}{
		{
			name:        "ValidYAML",
			fileContent: "name: Test\nage: 30",
			expectErr:   false,
			expected: map[string]interface{}{
				"name": "Test",
				"age":  30,
			},
		},
		{
			name:        "InvalidYAML",
			fileContent: "name: Test\nage: 30:",
			expectErr:   true,
		},
		{
			name:        "EmptyFile",
			fileContent: "",
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_yaml.yaml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.Write([]byte(tt.fileContent))
			if err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			values, err := valuesFromYamlFile(tmpFile.Name())
			if (err != nil) != tt.expectErr {
				t.Errorf("valuesFromYamlFile() error = %v, expectErr %v, file contect: %s", err, tt.expectErr, tt.fileContent)
				return
			}

			if tt.expectErr {
				return
			}

			if !equalMaps(values, tt.expected) {
				t.Errorf("valuesFromYamlFile() = %v, expected %v", values, tt.expected)
			}
		})
	}
}

func TestRemoveFilesFromDir(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(dir) // Clean up after the test

	// Create test files in the temporary directory
	files := []string{
		"docker-compose-rendered-1.yaml",
		"docker-compose-rendered-2.yaml",
		"docker-compose-rendered-3.yaml",
		"not-to-be-deleted.yaml",
	}
	for _, file := range files {
		_, err := os.Create(filepath.Join(dir, file))
		assert.NoError(t, err)
	}

	// Define the regex pattern to match files
	pattern := regexp.MustCompile(`^docker-compose-rendered-\d+\.yaml$`)

	// Call the function to remove matching files
	err = RemoveFilesFromDir(dir, pattern)
	assert.NoError(t, err)

	// Verify that the matching files were deleted
	for _, file := range files {
		_, err := os.Stat(filepath.Join(dir, file))
		if pattern.MatchString(file) {
			assert.True(t, os.IsNotExist(err), "File should have been deleted: %s", file)
		} else {
			assert.NoError(t, err, "File should not have been deleted: %s", file)
		}
	}
}

func TestRemoveFilesFromDir_NoMatches(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(dir) // Clean up after the test

	// Create test files in the temporary directory
	files := []string{
		"file1.yaml",
		"file2.yaml",
		"file3.yaml",
	}
	for _, file := range files {
		_, err := os.Create(filepath.Join(dir, file))
		assert.NoError(t, err)
	}

	// Define the regex pattern that doesn't match any files
	pattern := regexp.MustCompile(`^no-match-\d+\.yaml$`)

	// Call the function to remove matching files
	err = RemoveFilesFromDir(dir, pattern)
	assert.NoError(t, err)

	// Verify that no files were deleted
	for _, file := range files {
		_, err := os.Stat(filepath.Join(dir, file))
		assert.NoError(t, err, "File should not have been deleted: %s", file)
	}
}

func TestRemoveFilesFromDir_ErrorHandling(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(dir) // Clean up after the test

	// Define the regex pattern to match files
	pattern := regexp.MustCompile(`^docker-compose-rendered-\d+\.yaml$`)

	// Pass a non-existent directory to induce an error
	err = RemoveFilesFromDir(filepath.Join(dir, "non-existent-dir"), pattern)
	assert.Error(t, err, "Expected an error due to non-existent directory")
}

// Helper function to compare two maps
func equalMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
