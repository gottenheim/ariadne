package cmd

import (
	"github.com/gottenheim/ariadne/details/fs_repo"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var extractAnswerCmd = &cobra.Command{
	Use:   "extract-answer",
	Short: "Extracts files from archive to card directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{cardDirFlag})

		if err != nil {
			return err
		}

		cardDir := dirs[0]

		repo := fs_repo.NewFileCardRepository(osFs)

		section, entry := repo.GetCardPathSection(cardDir), repo.GetCardPathEntry(cardDir)

		useCase := &use_cases.ExtractCard{}

		return useCase.Run(repo, section, entry)
	},
}

func init() {
	cardCmd.AddCommand(extractAnswerCmd)

	extractAnswerCmd.Flags().String(cardDirFlag, "", "Card directory")
	extractAnswerCmd.MarkFlagRequired(cardDirFlag)
}
