package cmd

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func GetDirectoryFlag(cmd *cobra.Command, fs afero.Fs, flagName string) (string, error) {
	dirPath, _ := cmd.Flags().GetString(flagName)

	exists, err := afero.Exists(fs, dirPath)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("directory '%s' does not exist", dirPath)
	}

	return dirPath, nil
}

func GetDirectoryFlags(cmd *cobra.Command, fs afero.Fs, flagNames []string) ([]string, error) {
	var dirs []string

	for _, flagName := range flagNames {
		dirPath, err := GetDirectoryFlag(cmd, fs, flagName)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, dirPath)
	}

	return dirs, nil
}

func GetDirectoriesToIgnore(cmd *cobra.Command, flagName string) []string {
	dirsToIgnore, _ := cmd.Flags().GetStringArray(flagName)

	return dirsToIgnore
}
