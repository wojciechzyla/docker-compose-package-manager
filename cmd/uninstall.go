/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"log"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uninstallHelp = `
Uninstall docker compose project
`

func newUninstallCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "uninstall",
		Short: uninstallHelp,
		Long:  uninstallHelp,
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

				configurationFile = filepath.Join(packagePath, "rendered.yaml")
			}
			dockerCmd := exec.Command("docker", "compose", "-f", configurationFile, "down", "")
			output, err := dockerCmd.CombinedOutput()
			log.Printf("docker result:\n%s\n", output)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
		},
	}

	command.Flags().StringVarP(&composeConfig, "config", "c", "", "Path to the rendered configuration file")
	command.MarkFlagFilename("config")

	return command
}
