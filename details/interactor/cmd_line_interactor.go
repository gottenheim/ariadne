package interactor

import (
	"fmt"
	"io"
	"unicode"

	"github.com/eiannone/keyboard"
	sgr "github.com/foize/go.sgr"
	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
)

type CommandLineInteractor struct {
}

func NewCommandLineInteractor() *CommandLineInteractor {
	return &CommandLineInteractor{}
}

func (i *CommandLineInteractor) ShowDiscoveredDailyCards(dailyCards *study.DailyCards) {
	fmt.Printf(sgr.MustParseln("[fg-red]%d hot cards[reset], [fg-green]%d new cards[reset], [fg-blue]%d cards to remind[reset]"),
		len(dailyCards.HotCardsToRevise), len(dailyCards.NewCards), len(dailyCards.ScheduledCards))
}

func (i *CommandLineInteractor) AskQuestion(crd *card.Card, states []*study.CardState) (*study.CardState, error) {
	i.showQuestionHeader(crd)
	i.showQuestion(crd)
	err := i.waitForUserToComeUpWithAnswer()
	if err != nil {
		return nil, err
	}
	err = i.ShowAnswer(crd)
	if err != nil {
		return nil, err
	}
	return i.askUserHowGoodWasHisAnswer(crd, states)
}

func (i *CommandLineInteractor) ShowAnswer(crd *card.Card) error {
	answer, err := crd.Answer()
	if err != nil {
		return err
	}

	for name, content := range answer {
		fmt.Printf(sgr.MustParseln("[fg-green]----- %s -----\n"), name)
		fmt.Printf(sgr.MustParseln("%s\n\n"), content)
	}

	return nil
}

func (i *CommandLineInteractor) showQuestionHeader(crd *card.Card) {
	fmt.Println("------------------------------------------------")
	fmt.Printf(sgr.MustParseln("Question directory: [fg-white]%s/%s\n"), crd.Section(), crd.Entry())
}

func (i *CommandLineInteractor) showQuestion(crd *card.Card) {
	for _, artifact := range crd.Artifacts() {
		if artifact.Name() == card.AnswerArtifactName {
			continue
		}

		fmt.Printf(sgr.MustParseln("[fg-green]----- %s -----\n"), artifact.Name())
		fmt.Printf(sgr.MustParseln("%s\n\n"), artifact.Content())
	}
}

func (i *CommandLineInteractor) waitForUserToComeUpWithAnswer() error {
	fmt.Println(sgr.MustParseln("[fg-white]Press any key to show answer[reset]"))

	_, key, err := keyboard.GetSingleKey()
	if err != nil {
		return err
	}

	// ctrl + c or escape
	if key == 3 || key == 27 {
		return io.EOF
	}

	return nil
}

func (i *CommandLineInteractor) askUserHowGoodWasHisAnswer(crd *card.Card, states []*study.CardState) (*study.CardState, error) {
	gradeColors := map[int]string{
		0: "[fg-red]",
		3: "[fg-magenta]",
		4: "[fg-green]",
		5: "[fg-yellow]",
	}

	for _, state := range states {
		gradeColor, ok := gradeColors[state.Grade]
		if !ok {
			gradeColor = "[fg-white]"
		}
		fmt.Printf(sgr.MustParse(gradeColor+"[underline]%s[underlineOff]%s "), string(state.Name[0]), state.Name[1:])
	}

	fmt.Print("\n\n")

	for {
		ch, key, err := keyboard.GetSingleKey()
		if err != nil {
			return nil, err
		}

		// ctrl + c or escape
		if key == 3 || key == 27 {
			return nil, io.EOF
		}

		for _, state := range states {
			if state.Name[0] == byte(unicode.ToUpper(rune(ch))) {
				gradeColor, ok := gradeColors[state.Grade]
				if !ok {
					gradeColor = "[fg-white]"
				}

				fmt.Printf(sgr.MustParseln(gradeColor+"%s[reset] selected\n"), state.Name)
				return state, nil
			}
		}

		fmt.Println("Character was not recognized, try again")
	}
}
