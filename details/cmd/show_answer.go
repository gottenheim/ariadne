package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/fs/card_repo"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var showAnswerCmd = &cobra.Command{
	Use:   "show-answer",
	Short: "Shows card answer",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{"card-dir"})

		if err != nil {
			return err
		}

		cardDir := filepath.Dir(dirs[0])

		cardRepo := card_repo.NewFileCardRepository(osFs, cardDir)

		action := &card.ShowAnswerAction{}

		cardKey, err := strconv.Atoi(filepath.Base(dirs[0]))

		if err != nil {
			return errors.New("Card directory should be a number")
		}

		return action.Run(cardRepo, os.Stdout, card.Key(cardKey))
	},
}

func init() {
	cardCmd.AddCommand(showAnswerCmd)

	showAnswerCmd.Flags().String("card-dir", "", "Card directory")
	showAnswerCmd.MarkFlagRequired("card-dir")
}
