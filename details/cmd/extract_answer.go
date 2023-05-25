package cmd

import (
	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/fs/card_repo"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var extractAnswerCmd = &cobra.Command{
	Use:   "extract-answer",
	Short: "Extracts files from archive to card directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{"base-dir", "card-dir"})

		if err != nil {
			return err
		}

		baseDir, cardDir := dirs[0], dirs[1]

		cardRepo := card_repo.NewFileCardRepository(osFs, baseDir)

		action := &card.ExtractCardAction{}

		return action.Run(cardRepo, cardDir)
	},
}

func init() {
	cardCmd.AddCommand(extractAnswerCmd)

	extractAnswerCmd.Flags().String("base-dir", "", "Base directory (e.g. git repo directory)")
	extractAnswerCmd.MarkFlagRequired("base-dir")
	extractAnswerCmd.Flags().String("card-dir", "", "Card directory")
	extractAnswerCmd.MarkFlagRequired("card-dir")
}
