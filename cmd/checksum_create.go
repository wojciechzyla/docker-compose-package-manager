/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newChecksumCreateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create checksum of the package",
		Long:  "To create a checksum run \"dcpm checksum create [path_to_package]\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				return err
			}

			packagePath, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			if _, err := os.Stat(packagePath); errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("can't find a direcory: %s", packagePath)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var filesToSkip = []string{"CHECKSUM"}
			packagePath, err := filepath.Abs(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				os.Exit(1)
			}

			checksum, err := checksumDirectory(packagePath, filesToSkip)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				os.Exit(1)
			}

			file, err := os.Create(filepath.Join(packagePath, "CHECKSUM"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				os.Exit(1)
			}
			defer file.Close()

			_, err = file.WriteString(checksum)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err)
				os.Exit(1)
			}
		},
	}

	return command
}
