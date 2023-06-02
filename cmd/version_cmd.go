package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.0.2")
	},
}

func init() {
	rootCmd.AddCommand(showVersionCmd)
}
