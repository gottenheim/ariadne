package cmd

import (
	"github.com/gottenheim/ariadne/details/fs_repo"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const sectionDirFlag = "section-dir"
const templateDirFlag = "template-dir"

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new card",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{sectionDirFlag, templateDirFlag})

		if err != nil {
			return err
		}

		sectionDir, templateDir := dirs[0], dirs[1]

		useCase := &use_cases.NewCard{}

		templateRepo := fs_repo.NewFileTemplateRepository(osFs, templateDir)

		cardRepo := fs_repo.NewFileCardRepository(osFs)

		return useCase.Run(templateRepo, cardRepo, sectionDir)
	},
}

func init() {
	cardCmd.AddCommand(newCmd)

	newCmd.Flags().String(sectionDirFlag, "", "Cards section directory")
	newCmd.MarkFlagRequired(sectionDirFlag)
	newCmd.Flags().String(templateDirFlag, "", "Template files directory")
	newCmd.MarkFlagRequired(templateDirFlag)
}
