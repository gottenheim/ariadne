package cmd

import (
	"errors"
	"fmt"

	"github.com/gottenheim/ariadne/card"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var packAnswerCmd = &cobra.Command{
	Use:   "pack-answer",
	Short: "Packs files with correct answer to archive",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardDirPath, _ := cmd.Flags().GetString("card-dir")

		fs := afero.NewOsFs()

		exists, err := afero.Exists(fs, cardDirPath)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New(fmt.Sprintf("Directory '%s' does not exist", cardDirPath))
		}

		configFilePath, _ := cmd.Flags().GetString("config-file")

		config, err := card.LoadConfig(fs, configFilePath)
		if err != nil {
			return err
		}

		card := card.New(fs, config)

		err = card.PackAnswer(cardDirPath)
		if err != nil {
			return err
		}

		fmt.Printf("Card '%s' answer packed sucessfully\n", cardDirPath)

		return nil
	},
}

func init() {
	cardCmd.AddCommand(packAnswerCmd)

	packAnswerCmd.Flags().String("card-dir", "", "Card directory")
	packAnswerCmd.MarkFlagRequired("card-dir")
}
