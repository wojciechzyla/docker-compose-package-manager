/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newChecksumCheckCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "check",
		Short: "check checksum of the package",
		Long:  "To check a checksum run \"dcpm checksum check [path_to_package]\"",
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

			checksumPath := filepath.Join(packagePath, "CHECKSUM")
			if _, err := os.Stat(checksumPath); errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("can't find a CHECKSUM file: %s", checksumPath)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var filesToSkip = []string{"CHECKSUM"}
			packagePath, err := filepath.Abs(args[0])
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			checksum, err := checksumDirectory(packagePath, filesToSkip)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			file, err := os.Open(filepath.Join(packagePath, "CHECKSUM"))
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			if checksum == string(content) {
				log.Print("checksum is correct")
			} else {
				log.Fatalf("checksum doesn't match!")
			}
		},
	}

	return command
}
