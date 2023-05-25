package cmd

import (
	"strings"

	"github.com/gottenheim/ariadne/card"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new card",
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, fs, []string{"base-dir", "cards-dir", "template-dir"})

		if err != nil {
			return err
		}

		baseDir, cardsDir, templateDir := dirs[0], dirs[1], dirs[2]

		action := &card.NewCardAction{}

		templateRepo := card.NewFileTemplateRepository(fs, templateDir)

		cardRepo := card.NewFileCardRepository(fs, baseDir)

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
