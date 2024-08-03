/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dcpm",
	Short: "",
	Long:  `CLI aplication for managing docker compose packages`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		newInitCommand(),
		newRenderCommand(),
		newInstallCommand(),
		newUninstallCommand(),
		newChecksumCommand())
}
