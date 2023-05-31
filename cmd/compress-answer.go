package cmd

import (
	"path/filepath"

	"github.com/gottenheim/ariadne/infra/repo/fs_repo"
	"github.com/gottenheim/ariadne/use_cases"
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
		section := filepath.Dir(cardDir)
		entry := filepath.Base(cardDir)

		repo := fs_repo.NewFileCardRepository(osFs)

		useCase := &use_cases.CompressAnswer{}

		return useCase.Run(repo, section, entry)
	},
}

func init() {
	cardCmd.AddCommand(compressAnswerCmd)

	compressAnswerCmd.Flags().String("card-dir", "", "Card directory (relative to base)")
	compressAnswerCmd.MarkFlagRequired("card-dir")
}
