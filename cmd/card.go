package cmd

import (
	"github.com/spf13/cobra"
)

var cardCmd = &cobra.Command{
	Use: "card",
}

func init() {
	rootCmd.AddCommand(cardCmd)
}
