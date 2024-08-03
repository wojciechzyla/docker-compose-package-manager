/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func packagePathValid(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("direcotry doesn't exist: %s", path)
	}
	errorsSlice := make([]error, 0)
	if _, err := os.Stat(filepath.Join(path, "values.yaml")); errors.Is(err, os.ErrNotExist) {
		errorsSlice = append(errorsSlice, fmt.Errorf("can't find values.yaml inside direcotry: %s", path))
	}
	if _, err := os.Stat(filepath.Join(path, "templates")); errors.Is(err, os.ErrNotExist) {
		errorsSlice = append(errorsSlice, fmt.Errorf("can't find templates direcotry inside direcotry: %s", path))
	}
	if _, err := os.Stat(filepath.Join(path, "running_config")); errors.Is(err, os.ErrNotExist) {
		errorsSlice = append(errorsSlice, fmt.Errorf("can't find running_config direcotry inside direcotry: %s", path))
	}
	if _, err := os.Stat(filepath.Join(path, "dependencies")); errors.Is(err, os.ErrNotExist) {
		errorsSlice = append(errorsSlice, fmt.Errorf("can't find dependencies direcotry inside direcotry: %s", path))
	}
	if len(errorsSlice) > 0 {
		return errors.Join(errorsSlice...)
	}
	return nil
}

func filePathValid(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("can't find a file: %s", path)
	}
	return nil
}

func processFilePath(path *string) error {
	absPath, err := filepath.Abs(*path)
	if err != nil {
		return err
	}
	*path = absPath
	err = filePathValid(*path)
	if err != nil {
		return err
	}
	return nil
}
