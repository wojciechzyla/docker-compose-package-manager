package src

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
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

func TestCombineYamls(t *testing.T) {
	tests := []struct {
		name            string
		sourceFiles     map[string]string
		expectErr       bool
		expectedContent string
	}{
		{
			name: "CombineMultipleYamls",
			sourceFiles: map[string]string{
				"file1.yaml": "name: Test1\nage: 30",
				"file2.yaml": "name: Test2\nage: 25",
			},
			expectErr:       false,
			expectedContent: "name: Test1\nage: 30\n---\nname: Test2\nage: 25",
		},
		{
			name: "SingleYaml",
			sourceFiles: map[string]string{
				"file1.yaml": "name: Test1\nage: 30",
			},
			expectErr:       false,
			expectedContent: "name: Test1\nage: 30",
		},
		{
			name:            "NoYamls",
			sourceFiles:     map[string]string{},
			expectErr:       false,
			expectedContent: "",
		},
		{
			name:            "InvalidDirectory",
			sourceFiles:     nil,
			expectErr:       true,
			expectedContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceDir, err := os.MkdirTemp("", "test_src")
			if err != nil {
				t.Fatalf("Failed to create temp source directory: %v", err)
			}
			defer os.RemoveAll(sourceDir)

			for fileName, content := range tt.sourceFiles {
				filePath := filepath.Join(sourceDir, fileName)
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to write to source file: %v", err)
				}
			}

			destFile, err := os.CreateTemp("", "test_dest.yaml")
			if err != nil {
				t.Fatalf("Failed to create temp destination file: %v", err)
			}
			defer os.Remove(destFile.Name())
			destFile.Close()

			srcDir := sourceDir
			if tt.sourceFiles == nil {
				srcDir = "invalid_directory"
			}

			err = combineYamls(srcDir, destFile.Name())
			if (err != nil) != tt.expectErr {
				t.Errorf("combineYamls() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr {
				content, err := os.ReadFile(destFile.Name())
				if err != nil {
					t.Fatalf("Failed to read destination file: %v", err)
				}
				if strings.TrimSpace(string(content)) != tt.expectedContent {
					t.Errorf("combineYamls() content = %v, expected %v", string(content), tt.expectedContent)
				}
			}
		})
	}
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
