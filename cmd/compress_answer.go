package cmd

import (
	"github.com/gottenheim/ariadne/details/fs_repo"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var compressAnswerCmd = &cobra.Command{
	Use:   "compress-answer",
	Short: "Compresses card files to answer file",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{cardDirFlag})

		if err != nil {
			return err
		}

		repo := fs_repo.NewFileCardRepository(osFs)

		cardDir := dirs[0]

		section, entry := repo.GetCardPathSection(cardDir), repo.GetCardPathEntry(cardDir)

		useCase := &use_cases.CompressAnswer{}

		return useCase.Run(repo, section, entry)
	},
}

func init() {
	cardCmd.AddCommand(compressAnswerCmd)

	compressAnswerCmd.Flags().String(cardDirFlag, "", "Card directory (relative to base)")
	compressAnswerCmd.MarkFlagRequired(cardDirFlag)
}
