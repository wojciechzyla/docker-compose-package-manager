/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
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
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configurationDir := ""
			pattern := regexp.MustCompile(`^docker-compose-rendered-\d+\.yaml$`)
			if len(composeConfig) > 0 {
				err := processFilePath(&composeConfig)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v", err)
					os.Exit(1)
				}
				configurationDir = composeConfig
			} else {
				packagePath, err := filepath.Abs(args[0])
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v", err)
					os.Exit(1)
				}

				err = packagePathValid(packagePath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %v", err)
					os.Exit(1)
				}

				configurationDir = filepath.Join(packagePath, "running_config")
			}
			files, err := dockerComposeFilesInstall(configurationDir, pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: during handling compose files: %v", err)
				os.Exit(1)
			}
			files = append([]string{"compose"}, files...)
			files = append(files, "down")
			dockerCmd := exec.Command("docker", files...)
			output, err := dockerCmd.CombinedOutput()
			fmt.Printf("docker result:\n%s\n", output)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				os.Exit(1)
			}
			err = src.RemoveFilesFromDir(configurationDir, pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: during deleting files: %v", err)
				os.Exit(1)
			}
		},
	}

	command.Flags().StringVarP(&composeConfig, "config", "c", "", "Path to the rendered configuration file")
	command.MarkFlagFilename("config")

	return command
}
