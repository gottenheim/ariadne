package cmd

import (
	"os"
	"path/filepath"

	"github.com/gottenheim/ariadne/infra/fs/fs_repo"
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

		repo := fs_repo.NewFileCardRepository(osFs)

		section, entry := repo.GetCardPathSection(cardDir), repo.GetCardPathEntry(cardDir)

		useCase := &use_cases.ShowAnswer{}

		return useCase.Run(repo, os.Stdout, section, entry)
	},
}

func init() {
	cardCmd.AddCommand(showAnswerCmd)

	showAnswerCmd.Flags().String("card-dir", "", "Card directory")
	showAnswerCmd.MarkFlagRequired("card-dir")
}
