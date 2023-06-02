package cmd

import (
	"path/filepath"

	"github.com/gottenheim/ariadne/details/fs_repo"
	"github.com/gottenheim/ariadne/details/interactor"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const cardDirFlag = "card-dir"

var showAnswerCmd = &cobra.Command{
	Use:   "show-answer",
	Short: "Shows card answer",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{cardDirFlag})

		if err != nil {
			return err
		}

		cardDir := filepath.Dir(dirs[0])

		repo := fs_repo.NewFileCardRepository(osFs)

		section, entry := repo.GetCardPathSection(cardDir), repo.GetCardPathEntry(cardDir)

		userInteractor := interactor.NewCommandLineInteractor()

		useCase := &use_cases.ShowAnswer{}

		return useCase.Run(repo, userInteractor, section, entry)
	},
}

func init() {
	cardCmd.AddCommand(showAnswerCmd)

	showAnswerCmd.Flags().String(cardDirFlag, "", "Card directory")
	showAnswerCmd.MarkFlagRequired(cardDirFlag)
}
