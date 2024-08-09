/*
Copyright © 2024 Wojciech Żyła <wojciechzyla.mail@gmail.com>
*/
package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func checksumDirectory(root string, filesToSkip []string) (string, error) {
	hasher := sha256.New()
	err := filepath.WalkDir(root, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		for _, f := range filesToSkip {
			if info.Name() == f {
				return nil
			}
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if !info.IsDir() {
			if _, err := io.Copy(hasher, file); err != nil {
				return err
			}
		}

		if _, err := hasher.Write([]byte(path)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

var checksumHelp = `
To create a checksum run "dcpm checksum create [path_to_package]".

To verify a checksum run "dcpm checksum check [path_to_package]".
`

func newChecksumCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "checksum",
		Short: "Create or verify checksum of the package",
		Long:  checksumHelp,
		Args:  cobra.NoArgs,
	}
	command.AddCommand(
		newChecksumCreateCommand(),
		newChecksumCheckCommand(),
	)
	return command
}
