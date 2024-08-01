/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wojciechzyla/docker-compose-package-manager/src"
)

var renderHelp = `
Render configuration file and save it to the provided directory. 
If the ouput file doesn't exist, it will be created.
`

var packagePath string
var outputPath string
var customValues string

func newRenderCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "render",
		Short: "Render configuration file",
		Long:  renderHelp,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if !strings.HasPrefix(packagePath, "/") {
				cwd, err := os.Getwd()
				if err != nil {
					log.Fatalf("error: %v", err)
					return
				}
				packagePath = filepath.Join(cwd, packagePath)
			}
			if !packagePathValid(packagePath) {
				return
			}

			if len(customValues) > 0 {
				if !strings.HasPrefix(customValues, "/") {
					cwd, err := os.Getwd()
					if err != nil {
						log.Fatalf("error: %v", err)
						return
					}
					customValues = filepath.Join(cwd, customValues)
				}
				if _, err := os.Stat(customValues); os.IsNotExist(err) {
					log.Fatalf("can't find a file: %s", customValues)
					return
				}
			}

			err := src.Render(packagePath, outputPath, customValues)
			if err != nil {
				log.Fatalf("error occured while parsing files: %v", err)
			}
		},
	}

	command.Flags().StringVarP(&packagePath, "package_path", "p", "", "Path of the docker compose package (required)")
	command.MarkFlagRequired("package_path")
	command.MarkFlagDirname("package_path")

	command.Flags().StringVarP(&outputPath, "output_path", "o", "", "Path the output file (required)")
	command.MarkFlagRequired("output_path")
	command.MarkFlagFilename("output_path")

	command.Flags().StringVarP(&customValues, "values", "v", "", "Path to the custom values.yaml")
	command.MarkFlagFilename("values")

	return command
}

func packagePathValid(path string) bool {
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		fmt.Printf("direcotry doesn't exist: %s", path)
		return false
	}
	errors := make([]string, 0)
	result := true
	if _, err := os.Stat(filepath.Join(path, "values.yaml")); os.IsNotExist(err) {
		errors = append(errors, fmt.Sprintf("can't find values.yaml inside direcotry: %s", path))
		result = false
	}
	if _, err := os.Stat(filepath.Join(path, "templates")); os.IsNotExist(err) {
		errors = append(errors, fmt.Sprintf("can't find templates direcotry inside direcotry: %s", path))
		result = false
	}
	if _, err := os.Stat(filepath.Join(path, "running_config")); os.IsNotExist(err) {
		errors = append(errors, fmt.Sprintf("can't find running_config direcotry inside direcotry: %s", path))
		result = false
	}
	if _, err := os.Stat(filepath.Join(path, "dependencies")); os.IsNotExist(err) {
		errors = append(errors, fmt.Sprintf("can't find dependencies direcotry inside direcotry: %s", path))
		result = false
	}
	if !result {
		for _, err := range errors {
			log.Fatalf("error: %v", err)
		}
	}
	return result
}
