package cmd

import (
	"strings"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/fs/card_repo"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new card",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{"base-dir", "cards-dir", "template-dir"})

		if err != nil {
			return err
		}

		baseDir, cardsDir, templateDir := dirs[0], dirs[1], dirs[2]

		action := &card.NewCardAction{}

		templateRepo := card_repo.NewFileTemplateRepository(osFs, templateDir)

		cardRepo := card_repo.NewFileCardRepository(osFs, baseDir)

		return action.Run(templateRepo, cardRepo, strings.Split(cardsDir, afero.FilePathSeparator))
	},
}

func init() {
	cardCmd.AddCommand(newCmd)

	newCmd.Flags().String("base-dir", "", "Base directory (e.g. git repo directory)")
	newCmd.MarkFlagRequired("base-dir")
	newCmd.Flags().String("cards-dir", "", "Cards subdirectory relative to base")
	newCmd.MarkFlagRequired("cards-dir")
	newCmd.Flags().String("template-dir", "", "Template files directory")
	newCmd.MarkFlagRequired("template-dir")
}
