/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
			packagePath, err := filepath.Abs(packagePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				os.Exit(1)
			}

			err = packagePathValid(packagePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				os.Exit(1)
			}

			if len(customValues) > 0 {
				err := processFilePath(&customValues)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v", err)
					os.Exit(1)
				}
			}
			err = src.Render(packagePath, outputPath, customValues)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: during parsing files: %v", err)
				os.Exit(1)
			}
		},
	}

	command.Flags().StringVarP(&packagePath, "package_path", "p", "", "Path of the docker compose package (required)")
	command.MarkFlagRequired("package_path")
	command.MarkFlagDirname("package_path")

	command.Flags().StringVarP(&outputPath, "output_path", "o", "", "Path to the directory, where rendered files will be saved (required)")
	command.MarkFlagRequired("output_path")
	command.MarkFlagFilename("output_path")

	command.Flags().StringVarP(&customValues, "values", "v", "", "Path to the custom values.yaml")
	command.MarkFlagFilename("values")

	return command
}
