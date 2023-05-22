package cmd

import (
	"errors"
	"fmt"

	"github.com/gottenheim/ariadne/card"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new card",
	RunE: func(cmd *cobra.Command, args []string) error {
		cardsDirPath, _ := cmd.Flags().GetString("cards-dir")

		fs := afero.NewOsFs()

		exists, err := afero.Exists(fs, cardsDirPath)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New(fmt.Sprintf("Directory '%s' does not exist", cardsDirPath))
		}

		templateDirPath, _ := cmd.Flags().GetString("template-dir")

		exists, err = afero.Exists(fs, templateDirPath)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New(fmt.Sprintf("Directory '%s' does not exist", templateDirPath))
		}

		configFilePath, _ := cmd.Flags().GetString("config-file")

		config, err := card.LoadConfig(fs, configFilePath)
		if err != nil {
			return err
		}

		card := card.New(fs, config)

		cardDir, err := card.CreateCard(cardsDirPath, templateDirPath)

		if err != nil {
			return err
		}

		fmt.Printf("Card '%s' has been created\n", cardDir)

		return nil
	},
}

func init() {
	cardCmd.AddCommand(newCmd)

	newCmd.Flags().String("cards-dir", "", "Directory to create a new card")
	newCmd.MarkFlagRequired("cards-dir")
	newCmd.Flags().String("template-dir", "", "Directory to copy template from")
	newCmd.MarkFlagRequired("template-dir")
}
