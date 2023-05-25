package cmd

import (
	"os"

	"github.com/gottenheim/ariadne/card"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var showAnswerCmd = &cobra.Command{
	Use:   "show-answer",
	Short: "Shows card answer",
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, fs, []string{"base-dir", "card-dir"})

		if err != nil {
			return err
		}

		baseDir, cardDir := dirs[0], dirs[1]

		cardRepo := card.NewFileCardRepository(fs, baseDir)

		action := &card.ShowAnswerAction{}

		return action.Run(cardRepo, os.Stdout, cardDir)
	},
}

func init() {
	cardCmd.AddCommand(showAnswerCmd)

	showAnswerCmd.Flags().String("base-dir", "", "Base directory (e.g. git repo directory)")
	showAnswerCmd.MarkFlagRequired("base-dir")
	showAnswerCmd.Flags().String("card-dir", "", "Card directory")
	showAnswerCmd.MarkFlagRequired("card-dir")
}
