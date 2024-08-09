/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/wojciechzyla/docker-compose-package-manager/src"
)

var installHelp = `Install docker compose project`
var composeConfig string
var installValues string

func newInstallCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "install",
		Short: installHelp,
		Long:  installHelp,
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

				if len(installValues) > 0 {
					err := processFilePath(&installValues)
					if err != nil {
						log.Fatalf("error: %v", err)
					}
				}
				configurationDir = filepath.Join(packagePath, "running_config")
				err = src.RemoveFilesFromDir(configurationDir, pattern)
				if err != nil {
					log.Fatalf("error occured while deleting files: %v", err)
				}
				err = src.Render(packagePath, configurationDir, installValues)
				if err != nil {
					log.Fatalf("error occured while parsing files: %v", err)
				}
			}
			files, err := dockerComposeFilesInstall(configurationDir, pattern)
			if err != nil {
				log.Fatalf("error: during handling compose files: %v", err)
			}
			files = append([]string{"compose"}, files...)
			files = append(files, "up", "-d")
			dockerCmd := exec.Command("docker", files...)
			output, err := dockerCmd.CombinedOutput()
			log.Printf("docker result:\n%s\n", output)
			if err != nil {
				log.Printf("error: %v", err)
				err := src.RemoveFilesFromDir(configurationDir, pattern)
				if err != nil {
					log.Fatalf("error occured while deleting files: %v", err)
				}
				return
			}
		},
	}

	command.Flags().StringVarP(&composeConfig, "config", "c", "", "Path to the direcotry where configuration files will be rendered")
	command.MarkFlagFilename("config")

	command.Flags().StringVarP(&installValues, "values", "v", "", "Path to the values.yaml")
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
