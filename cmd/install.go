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

var installHelp = `To install docker compose project run 'dcpm install [path_to_package]'`
var composeConfig string
var installValues string

func newInstallCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "install",
		Short: `Install docker compose project`,
		Long:  installHelp,
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

				if len(installValues) > 0 {
					err := processFilePath(&installValues)
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %v", err)
						os.Exit(1)
					}
				}
				configurationDir = filepath.Join(packagePath, "running_config")
				err = src.RemoveFilesFromDir(configurationDir, pattern)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: during deleting files: %v", err)
					os.Exit(1)
				}
				err = src.Render(packagePath, configurationDir, installValues)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: during parsing files: %v", err)
					os.Exit(1)
				}
			}
			files, err := dockerComposeFilesInstall(configurationDir, pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: during handling compose files: %v", err)
				os.Exit(1)
			}
			files = append([]string{"compose"}, files...)
			files = append(files, "up", "-d")
			dockerCmd := exec.Command("docker", files...)
			output, err := dockerCmd.CombinedOutput()
			fmt.Printf("docker result:\n%s\n", output)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				err := src.RemoveFilesFromDir(configurationDir, pattern)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: during deleting files: %v", err)
					os.Exit(1)
				}
				os.Exit(1)
			}
		},
	}

	command.Flags().StringVarP(&composeConfig, "config", "c", "", "Path to the direcotry where configuration files will be rendered (optional)")
	command.MarkFlagFilename("config")

	command.Flags().StringVarP(&installValues, "values", "v", "", "Path to the cutom values.yaml (optional)")
	command.MarkFlagFilename("values")

	return command
}

func dockerComposeFilesInstall(directory string, pattern *regexp.Regexp) ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(directory, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && pattern.MatchString(info.Name()) {
			files = append(files, "-f", path)
		}
		return nil
	})
	return files, err
}
