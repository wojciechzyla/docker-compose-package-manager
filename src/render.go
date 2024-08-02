package src

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"dario.cat/mergo"
	"github.com/pkg/errors"
)

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
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	for i, template := range templates {
		fname := filepath.Join(tmpDir, fmt.Sprintf("rendered%d.yaml", i))
		if err != nil {
			return errors.Wrap(err, "creating output tmp file")
		}
		err = executeTemplate(fname, values, template)
		if err != nil {
			return errors.Wrap(err, "executing template file")
		}
	}
	err = combineYamls(tmpDir, destinationPath)
	if err != nil {
		return err
	}
	return nil
}

func executeTemplate(filePath string, values map[string]interface{}, tmpl *template.Template) error {
	var tmpOutput bytes.Buffer
	err := tmpl.Execute(&tmpOutput, values)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(tmpOutput.String())) > 0 {
		err := os.WriteFile(filePath, tmpOutput.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
