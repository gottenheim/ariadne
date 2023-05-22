package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/gottenheim/ariadne/card"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var unpackAnswerCmd = &cobra.Command{
	Use:   "unpack-answer",
	Short: "Unpacks files from archive to card directory",
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

		ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

		card := card.New(fs, config, ioStreams)

		err = card.UnpackAnswer(cardDirPath)
		if err != nil {
			return err
		}

		fmt.Printf("Card '%s' answer unpacked sucessfully\n", cardDirPath)

		return nil
	},
}

func init() {
	cardCmd.AddCommand(unpackAnswerCmd)

	unpackAnswerCmd.Flags().String("card-dir", "", "Card directory")
	unpackAnswerCmd.MarkFlagRequired("card-dir")
}
