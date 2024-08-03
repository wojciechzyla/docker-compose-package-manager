/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wojciechzyla/docker-compose-package-manager/src"
)

var installHelp = `
Install docker compose project
`
var composeConfig string
var installValues string

func newInstallCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "install",
		Short: installHelp,
		Long:  installHelp,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configurationFile := ""
			if len(composeConfig) > 0 {
				err := processFilePath(&composeConfig)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				configurationFile = composeConfig
			} else {
				if len(args) == 0 {
					log.Fatalf("no path to the package provided")
				}
				packagePath, err := filepath.Abs(args[0])
				if err != nil {
					log.Fatalf("error: %v", err)
				}

				err = packagePathValid(packagePath)
				if err != nil {
					log.Fatalf("error: %v", err)
				}

				if len(installValues) > 0 {
					err := processFilePath(&installValues)
					if err != nil {
						log.Fatalf("error: %v", err)
					}
				}
				configurationFile = filepath.Join(packagePath, "rendered.yaml")
				err = src.Render(packagePath, configurationFile, installValues)
				if err != nil {
					log.Fatalf("error occured while parsing files: %v", err)
				}
			}
			dockerCmd := exec.Command("docker", "compose", "-f", configurationFile, "up", "-d")
			output, err := dockerCmd.CombinedOutput()
			log.Printf("docker result:\n%s\n", output)
			if err != nil {
				log.Printf("error: %v", err)
				err := os.Remove(configurationFile)
				if err != nil {
					log.Fatalf("error: failed to remove file: %v", err)
				}
				return
			}
		},
	}

	command.Flags().StringVarP(&composeConfig, "config", "c", "", "Path to the rendered configuration file")
	command.MarkFlagFilename("config")

	command.Flags().StringVarP(&installValues, "values", "v", "", "Path to the values.yaml")
	command.MarkFlagFilename("values")

	return command
}
