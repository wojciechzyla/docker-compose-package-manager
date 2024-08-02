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

func traverseTemplatesDir(templatesPath string, templates *[]*template.Template) error {
	loop := true
	helperFiles := make([]string, 0)
	templateFiles := make([]string, 0)
	dirs := make([]string, 1)
	dirs[0] = templatesPath
	for loop {
		currentDir := dirs[len(dirs)-1]
		dirs = dirs[:len(dirs)-1]
		files, err := os.ReadDir(currentDir)
		if err != nil {
			return err
		}
		for _, file := range files {
			absPath, err := filepath.Abs(filepath.Join(currentDir, file.Name()))
			if err != nil {
				return err
			}
			if !file.IsDir() {
				if strings.HasSuffix(file.Name(), ".helper") {
					helperFiles = append(helperFiles, absPath)
				} else {
					templateFiles = append(templateFiles, absPath)
				}
			} else {
				dirs = append(dirs, absPath)
			}
		}
		files = nil
		if len(dirs) == 0 {
			loop = false
		}
	}
	err := parseFiles(templates, templateFiles, helperFiles)
	if err != nil {
		return err
	}
	return nil
}

func parseTemplates(templatesPath string) ([]*template.Template, error) {
	var templates []*template.Template
	err := traverseTemplatesDir(templatesPath, &templates)
	if err != nil {
		return nil, err
	}
	return templates, nil
}
