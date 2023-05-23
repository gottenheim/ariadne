package cmd

import (
	"github.com/gottenheim/ariadne/card"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var compressAnswerCmd = &cobra.Command{
	Use:   "compress-answer",
	Short: "Compresses card files to answer file",
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, fs, []string{"base-dir", "card-dir"})

		if err != nil {
			return err
		}

		action := &card.CompressAnswerAction{}

		return action.Run(fs, dirs[0], dirs[1])
	},
}

func init() {
	cardCmd.AddCommand(compressAnswerCmd)

	compressAnswerCmd.Flags().String("base-dir", "", "Base directory (e.g. git repo directory)")
	compressAnswerCmd.MarkFlagRequired("base-dir")
	compressAnswerCmd.Flags().String("card-dir", "", "Card directory (relative to base)")
	compressAnswerCmd.MarkFlagRequired("card-dir")
}
