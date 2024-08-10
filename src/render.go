package src

import (
	"fmt"
	"path/filepath"

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

	for i, template := range templates {
		fname := filepath.Join(destinationPath, fmt.Sprintf("docker-compose-rendered-%d.yaml", i))
		err = executeTemplate(fname, values, template)
		if err != nil {
			return errors.Wrap(err, "executing template file")
		}
	}
	return nil
}
