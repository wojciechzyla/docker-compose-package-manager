package src

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func executeTemplate(filePath string, values map[string]interface{}, tmpl *template.Template) error {
	var tmpOutput bytes.Buffer
	if values == nil {
		return errors.Errorf("providel values are nil")
	}
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

func RemoveFilesFromDir(directory string, pattern *regexp.Regexp) error {
	err := filepath.WalkDir(directory, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && pattern.MatchString(info.Name()) {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return errors.Wrap(err, "Error walking through directory:")
}
