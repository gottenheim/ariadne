package cmd

import (
	"path/filepath"

	"github.com/gottenheim/ariadne/infra/repo/fs_repo"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var extractAnswerCmd = &cobra.Command{
	Use:   "extract-answer",
	Short: "Extracts files from archive to card directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{"card-dir"})

		if err != nil {
			return err
		}

		cardDir := filepath.Dir(dirs[0])
		section := filepath.Dir(cardDir)
		entry := filepath.Base(cardDir)

		cardRepo := fs_repo.NewFileCardRepository(osFs)

		useCase := &use_cases.ExtractCard{}

		return useCase.Run(cardRepo, section, entry)
	},
}

func init() {
	cardCmd.AddCommand(extractAnswerCmd)

	extractAnswerCmd.Flags().String("card-dir", "", "Card directory")
	extractAnswerCmd.MarkFlagRequired("card-dir")
}
