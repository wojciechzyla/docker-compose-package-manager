/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"log"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/wojciechzyla/docker-compose-package-manager/src"
)

var uninstallHelp = `Uninstall docker compose project`

func newUninstallCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "uninstall",
		Short: uninstallHelp,
		Long:  uninstallHelp,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configurationDir := ""
			pattern := regexp.MustCompile(`^docker-compose-rendered-\d+\.yaml$`)
			if len(composeConfig) > 0 {
				err := processFilePath(&composeConfig)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				configurationDir = composeConfig
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

				configurationDir = filepath.Join(packagePath, "running_config")
			}
			files, err := dockerComposeFilesInstall(configurationDir, pattern)
			if err != nil {
				log.Fatalf("error: during handling compose files: %v", err)
			}
			files = append([]string{"compose"}, files...)
			files = append(files, "down")
			dockerCmd := exec.Command("docker", files...)
			output, err := dockerCmd.CombinedOutput()
			log.Printf("docker result:\n%s\n", output)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			err = src.RemoveFilesFromDir(configurationDir, pattern)
			if err != nil {
				log.Fatalf("error occured while deleting files: %v", err)
			}
		},
	}

	command.Flags().StringVarP(&composeConfig, "config", "c", "", "Path to the rendered configuration file")
	command.MarkFlagFilename("config")

	return command
}
