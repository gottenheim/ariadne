package cmd

import (
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

		action := &card.NewCardAction{}

		return action.Run(fs, dirs[0], dirs[1], dirs[2])
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
