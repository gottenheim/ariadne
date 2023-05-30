package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/infra/repo/fs_repo"
	"github.com/gottenheim/ariadne/use_cases"
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

		cardRepo := fs_repo.NewFileCardRepository(osFs, cardDir)

		action := &use_cases.ShowAnswerAction{}

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
