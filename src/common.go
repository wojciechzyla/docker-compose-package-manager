package src

import (
	"io"
	"os"
	"path/filepath"
	"strings"

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

func combineYamls(sourceDirectory string, destinationFilePath string) error {
	files, err := os.ReadDir(sourceDirectory)
	if err != nil {
		return err
	}
	var combinedContent strings.Builder

	for i, file := range files {
		if !file.IsDir() {
			absPath, err := filepath.Abs(filepath.Join(sourceDirectory, file.Name()))
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
	err = os.WriteFile(destinationFilePath, []byte(combinedContent.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}
