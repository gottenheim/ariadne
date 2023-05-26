package cmd

import (
	"errors"
	"path/filepath"
	"strconv"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/fs/card_repo"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var compressAnswerCmd = &cobra.Command{
	Use:   "compress-answer",
	Short: "Compresses card files to answer file",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{"card-dir"})

		if err != nil {
			return err
		}

		cardDir := filepath.Dir(dirs[0])

		repo := card_repo.NewFileCardRepository(osFs, cardDir)

		action := &card.CompressAnswerAction{}

		cardKey, err := strconv.Atoi(filepath.Base(dirs[0]))

		if err != nil {
			return errors.New("Card directory should be a number")
		}

		return action.Run(repo, card.Key(cardKey))
	},
}

func init() {
	cardCmd.AddCommand(compressAnswerCmd)

	compressAnswerCmd.Flags().String("card-dir", "", "Card directory (relative to base)")
	compressAnswerCmd.MarkFlagRequired("card-dir")
}
