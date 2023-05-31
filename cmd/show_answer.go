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

		section := filepath.Dir(cardDir)
		entry := filepath.Base(cardDir)

		cardRepo := fs_repo.NewFileCardRepository(osFs)

		useCase := &use_cases.ShowAnswer{}

		return useCase.Run(cardRepo, os.Stdout, section, entry)
	},
}

func init() {
	cardCmd.AddCommand(showAnswerCmd)

	showAnswerCmd.Flags().String("card-dir", "", "Card directory")
	showAnswerCmd.MarkFlagRequired("card-dir")
}
