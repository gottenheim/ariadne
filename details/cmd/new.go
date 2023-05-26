package cmd

import (
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

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{"cards-dir", "template-dir"})

		if err != nil {
			return err
		}

		cardsDir, templateDir := dirs[0], dirs[1]

		action := &card.NewCardAction{}

		templateRepo := card_repo.NewFileTemplateRepository(osFs, templateDir)

		cardRepo := card_repo.NewFileCardRepository(osFs, cardsDir)

		return action.Run(templateRepo, cardRepo)
	},
}

func init() {
	cardCmd.AddCommand(newCmd)

	newCmd.Flags().String("cards-dir", "", "Cards subdirectory relative to base")
	newCmd.MarkFlagRequired("cards-dir")
	newCmd.Flags().String("template-dir", "", "Template files directory")
	newCmd.MarkFlagRequired("template-dir")
}
