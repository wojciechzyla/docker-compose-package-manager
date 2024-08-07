package src

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"
)

func TestParseFiles(t *testing.T) {
	tests := []struct {
		name             string
		templateFiles    map[string]string
		helperFiles      map[string]string
		expectErr        bool
		expectedTplCount int
	}{
		{
			name: "MultipleTemplatesWithHelpers",
			templateFiles: map[string]string{
				"tmpl1.html": "{{define \"tmpl1\"}}Template 1 uses {{template \"helper\" .}}{{end}}",
				"tmpl2.html": "{{define \"tmpl2\"}}Template 2 uses {{template \"helper\" .}}{{end}}",
			},
			helperFiles: map[string]string{
				"helper1.html": "{{define \"helper\"}}Helper 1{{end}}",
			},
			expectErr:        false,
			expectedTplCount: 2,
		},
		{
			name: "SingleTemplateWithoutHelpers",
			templateFiles: map[string]string{
				"tmpl1.html": "{{define \"tmpl1\"}}Template 1{{end}}",
			},
			helperFiles:      map[string]string{},
			expectErr:        false,
			expectedTplCount: 1,
		},
		{
			name: "TemplateParseError",
			templateFiles: map[string]string{
				"tmpl1.html": "{{define \"tmpl1\"}Template 1{{end}}", // Missing closing bracket
			},
			helperFiles:      map[string]string{},
			expectErr:        true,
			expectedTplCount: 0,
		},
		{
			name: "HelperParseError",
			templateFiles: map[string]string{
				"tmpl1.html": "{{define \"tmpl1\"}}Template 1 uses {{template \"helper\" .}}{{end}}",
			},
			helperFiles: map[string]string{
				"helper1.html": "{{define \"helper\"}Helper 1{{end}}", // Missing closing bracket
			},
			expectErr:        true,
			expectedTplCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for template and helper files
			tmpDir, err := os.MkdirTemp("", "test_templates")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Write template files
			for fileName, content := range tt.templateFiles {
				filePath := filepath.Join(tmpDir, fileName)
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to write template file: %v", err)
				}
			}

			// Write helper files
			for fileName, content := range tt.helperFiles {
				filePath := filepath.Join(tmpDir, fileName)
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to write helper file: %v", err)
				}
			}

			// Collect template and helper file paths
			templateFilePaths := make([]string, 0, len(tt.templateFiles))
			for fileName := range tt.templateFiles {
				templateFilePaths = append(templateFilePaths, filepath.Join(tmpDir, fileName))
			}
			helperFilePaths := make([]string, 0, len(tt.helperFiles))
			for fileName := range tt.helperFiles {
				helperFilePaths = append(helperFilePaths, filepath.Join(tmpDir, fileName))
			}

			// Call the function
			var templates []*template.Template
			err = parseFiles(&templates, templateFilePaths, helperFilePaths)
			if (err != nil) != tt.expectErr {
				t.Errorf("parseFiles() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			// Check the number of parsed templates
			if len(templates) != tt.expectedTplCount {
				t.Errorf("parseFiles() parsed %d templates, expected %d", len(templates), tt.expectedTplCount)
			}
		})
	}
}

func TestParseTemplates(t *testing.T) {
	tests := []struct {
		name             string
		files            map[string]string
		expectErr        bool
		expectedTplCount int
	}{
		{
			name: "NestedDirectoriesWithTemplatesAndHelpers",
			files: map[string]string{
				// Level 1
				"templates/tmpl1.html":     "{{define \"tmpl1\"}}Template 1 uses {{template \"helper\" .}}{{end}}",
				"templates/helper1.helper": "{{define \"helper\"}}Helper 1{{end}}",

				// Level 2
				"templates/dir1/tmpl2.html":     "{{define \"tmpl2\"}}Template 2 uses {{template \"helper\" .}}{{end}}",
				"templates/dir1/helper2.helper": "{{define \"helper\"}}Helper 2{{end}}",

				// Level 3
				"templates/dir1/dir2/tmpl3.html":     "{{define \"tmpl3\"}}Template 3 uses {{template \"helper\" .}}{{end}}",
				"templates/dir1/dir2/helper3.helper": "{{define \"helper\"}}Helper 3{{end}}",
			},
			expectErr:        false,
			expectedTplCount: 3,
		},
		{
			name:             "EmptyDirectory",
			files:            map[string]string{}, // No files
			expectErr:        false,
			expectedTplCount: 0,
		},
		{
			name: "TemplateParseError",
			files: map[string]string{
				"templates/tmpl1.html": "{{define \"tmpl1\"}Template 1{{end}}", // Invalid syntax
			},
			expectErr:        true,
			expectedTplCount: 0,
		},
		{
			name: "HelperParseError",
			files: map[string]string{
				"templates/tmpl1.html":     "{{define \"tmpl1\"}}Template 1 uses {{template \"helper\" .}}{{end}}",
				"templates/helper1.helper": "{{define \"helper\"}Helper 1{{end}}", // Invalid syntax
			},
			expectErr:        true,
			expectedTplCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for the test files
			tmpDir, err := os.MkdirTemp("", "test_templates")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Write the test files
			for fileName, content := range tt.files {
				filePath := filepath.Join(tmpDir, fileName)
				if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
					t.Fatalf("Failed to create directories for %s: %v", filePath, err)
				}
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", filePath, err)
				}
			}

			// Call the function
			templates, err := parseTemplates(tmpDir)
			if (err != nil) != tt.expectErr {
				t.Errorf("parseTemplates() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			// Check the number of parsed templates
			if len(templates) != tt.expectedTplCount {
				t.Errorf("parseTemplates() parsed %d templates, expected %d", len(templates), tt.expectedTplCount)
			}
		})
	}
}
