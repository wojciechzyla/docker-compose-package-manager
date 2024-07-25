package parsing

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func valuesFromYamlFile(filePath string) (map[string]interface{}, error) {
	data, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}
	defer data.Close()
	s, err := io.ReadAll(data)
	if err != nil {
		return nil, errors.Wrap(err, "reading data file")
	}
	var values map[string]interface{}
	err = yaml.Unmarshal(s, &values)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling yaml file")
	}
	return values, nil
}

func parseTemplates(templatesPath string) ([]*template.Template, error) {
	files, err := os.ReadDir(templatesPath)
	if err != nil {
		return nil, err
	}
	var templates []*template.Template

	for _, file := range files {
		if !file.IsDir() {
			absPath, err := filepath.Abs(filepath.Join(templatesPath, file.Name()))
			if err != nil {
				return nil, err
			}
			tmpl, err := template.ParseFiles(absPath)
			if err != nil {
				return nil, err
			}
			templates = append(templates, tmpl)
		}
	}
	return templates, nil
}

func combineYamls(tmpDir string, destination string) error {
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return err
	}
	var combinedContent strings.Builder

	for i, file := range files {
		if !file.IsDir() {
			absPath, err := filepath.Abs(filepath.Join(tmpDir, file.Name()))
			if err != nil {
				return err
			}
			content, err := os.ReadFile(absPath)
			if err != nil {
				return err
			}

			combinedContent.Write(content)

			// Add YAML document separator if it's not the last file
			if i < len(files)-1 {
				combinedContent.WriteString("\n---\n")
			}
		}
	}
	err = os.WriteFile(destination, []byte(combinedContent.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}

func Render(packagePath, destinationPath, customValuesPath string) error {
	values, err := valuesFromYamlFile(filepath.Join(packagePath, "values.yaml"))
	if err != nil {
		return errors.Wrap(err, "reading default values.yaml")
	}

	if len(customValuesPath) > 0 {
		customValues, err := valuesFromYamlFile(customValuesPath)
		if err != nil {
			return errors.Wrap(err, "reading custom values.yaml")
		}
		if err := mergo.Merge(&values, customValues, mergo.WithOverride); err != nil {
			return errors.Wrap(err, "overwriting default values.yaml")
		}
	}

	templatesPath := filepath.Join(packagePath, "templates")
	templates, err := parseTemplates(templatesPath)
	if err != nil {
		return errors.Wrap(err, "paring emplates")
	}

	tmpDir, err := os.MkdirTemp("", "compose_render")
	check(err)
	defer os.RemoveAll(tmpDir)

	for i, template := range templates {
		fname := filepath.Join(tmpDir, fmt.Sprintf("rendered%d.yaml", i))
		output, err := os.Create(fname)
		if err != nil {
			return errors.Wrap(err, "creating output tmp file")
		}
		err = template.Execute(output, values)
		if err != nil {
			return errors.Wrap(err, "executing template file")
		}
		output.Close()
	}
	err = combineYamls(tmpDir, destinationPath)
	if err != nil {
		return err
	}
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
