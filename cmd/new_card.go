package cmd

import (
	"github.com/gottenheim/ariadne/infra/repo/fs_repo"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new card",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{"section-dir", "template-dir"})

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

	newCmd.Flags().String("section-dir", "", "Cards section directory")
	newCmd.MarkFlagRequired("section-dir")
	newCmd.Flags().String("template-dir", "", "Template files directory")
	newCmd.MarkFlagRequired("template-dir")
}
