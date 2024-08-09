/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initHelp = `Initialization of dcpm project with all basic directories and files. Run dcpm <project_name>`

func handleCreationError(message string, path string) {
	os.RemoveAll(path)
	log.Fatalf(message)
}

func newInitCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "init",
		Short: "Initialization of dcpm project",
		Long:  initHelp,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			newProjectPath := filepath.Join(cwd, args[0])
			if _, err := os.Stat(newProjectPath); !os.IsNotExist(err) {
				return fmt.Errorf("directory already exists: %s", args[0])
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			projectPath := filepath.Join(cwd, args[0])
			templatesPath := filepath.Join(projectPath, "templates")
			dependenciesPath := filepath.Join(projectPath, "additional")
			runningConfigPath := filepath.Join(projectPath, "running_config")

			if err := os.MkdirAll(templatesPath, os.ModePerm); err != nil {
				handleCreationError(fmt.Sprintf("error while creating a new directory: %s", err), projectPath)
			}

			if err := os.MkdirAll(dependenciesPath, os.ModePerm); err != nil {
				handleCreationError(fmt.Sprintf("error while creating a new directory: %s", err), projectPath)
			}

			if err := os.MkdirAll(runningConfigPath, os.ModePerm); err != nil {
				handleCreationError(fmt.Sprintf("error while creating a new directory: %s", err), projectPath)
			}

			valuesPath := filepath.Join(projectPath, "values.yaml")
			if _, err := os.Create(valuesPath); err != nil {
				handleCreationError(fmt.Sprintf("error while creating a new file: %s", err), projectPath)
			}
		},
	}

	return command
}
