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
			file := ""
			if len(composeConfig) > 0 {
				err := processFilePath(&composeConfig)
				if err != nil {
					log.Fatalf("error: %v", err)
					return
				}
				file = composeConfig
			} else {
				if len(args) == 0 {
					log.Fatalf("no path to the package provided")
					return
				}
				packagePath, err := filepath.Abs(args[0])
				if err != nil {
					log.Fatalf("error: %v", err)
					return
				}

				err = packagePathValid(packagePath)
				if err != nil {
					log.Fatalf("error: %v", err)
					return
				}

				if len(installValues) > 0 {
					err := processFilePath(&installValues)
					if err != nil {
						log.Fatalf("error: %v", err)
						return
					}
				}
				file = filepath.Join(packagePath, "rendered.yaml")
				err = src.Render(packagePath, file, installValues)
				if err != nil {
					log.Fatalf("error occured while parsing files: %v", err)
					return
				}
			}
			dockerCmd := exec.Command("docker", "compose", "-f", file, "up", "-d")
			output, err := dockerCmd.CombinedOutput()
			log.Printf("docker result:\n%s\n", output)
			if err != nil {
				log.Printf("error: %v", err)
				err := os.Remove(file)
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
