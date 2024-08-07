package src

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func parseFiles(templates *[]*template.Template, templateFiles, helperFiles []string) error {
	for _, filePath := range templateFiles {
		tmpFiles := make([]string, 1)
		tmpFiles[0] = filePath
		tmpFiles = append(tmpFiles, helperFiles...)
		tmpl, err := template.ParseFiles(tmpFiles...)
		if err != nil {
			return err
		}
		*templates = append(*templates, tmpl)
	}
	return nil
}

func parseTemplates(templatesPath string) ([]*template.Template, error) {
	var templates []*template.Template
	helperFiles := make([]string, 0)
	templateFiles := make([]string, 0)
	err := filepath.WalkDir(templatesPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if strings.HasSuffix(d.Name(), ".helper") {
			helperFiles = append(helperFiles, absPath)
		} else {
			templateFiles = append(templateFiles, absPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = parseFiles(&templates, templateFiles, helperFiles)
	if err != nil {
		return nil, err
	}
	return templates, nil
}
