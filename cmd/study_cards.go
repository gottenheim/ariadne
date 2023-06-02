package cmd

import (
	"fmt"

	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/details/fs_repo"
	"github.com/gottenheim/ariadne/details/interactor"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/use_cases"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const dirFlag = "dir"
const newCardsFlag = "new-cards"
const cardsToRemindFlag = "cards-to-remind"

var studyCardsCmd = &cobra.Command{
	Use:   "study-cards",
	Short: "Executes study cards session in directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		osFs := afero.NewOsFs()

		dirs, err := GetDirectoryFlags(cmd, osFs, []string{dirFlag})

		if err != nil {
			return err
		}

		cardsDir := dirs[0]

		fmt.Printf("Directory used to discover cards: %s\n", cardsDir)

		cardRepo := fs_repo.NewFileCardRepository(osFs)
		timeSource := datetime.NewOsTimeSource()

		useCase := use_cases.NewStudyCardsSession(timeSource, cardRepo, interactor.NewCommandLineInteractor())

		cardEmitter := fs_repo.NewAnsweredCardEmitter(osFs, cardRepo, cardsDir)
		newCards, _ := cmd.Flags().GetInt(newCardsFlag)
		cardsToRemind, _ := cmd.Flags().GetInt(cardsToRemindFlag)

		return useCase.Run(cardEmitter, &study.DailyCardsConfig{
			NewCardsCount:       newCards,
			ScheduledCardsCount: cardsToRemind,
		})
	},
}

func init() {
	rootCmd.AddCommand(studyCardsCmd)

	studyCardsCmd.Flags().String(dirFlag, "", "Card directory")
	studyCardsCmd.MarkFlagRequired(dirFlag)

	studyCardsCmd.Flags().Int(newCardsFlag, 0, "Count of new cards to study today")
	studyCardsCmd.MarkFlagRequired(newCardsFlag)

	studyCardsCmd.Flags().Int(cardsToRemindFlag, 0, "Count of cards to remind today")
	studyCardsCmd.MarkFlagRequired(cardsToRemindFlag)
}
