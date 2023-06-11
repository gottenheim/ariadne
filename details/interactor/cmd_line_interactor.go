package interactor

import (
	"fmt"
	"io"
	"strconv"
	"unicode"

	"github.com/atotto/clipboard"
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

func (i *CommandLineInteractor) ShowStudyProgress(selectedDailyCard *study.SelectedDailyCard, studyProgress *study.StudyProgress) {
	underlineIf := func(condition bool) string {
		if condition {
			return "[underline]"
		} else {
			return ""
		}
	}
	incrementIf := func(condition bool, val int) int {
		if condition {
			return val + 1
		} else {
			return val
		}
	}

	str := fmt.Sprintf("[fg-red]%s%d hot cards[reset], [fg-green]%s%d new cards[reset], [fg-blue]%s%d cards to remind[reset]",
		underlineIf(selectedDailyCard.CardType == study.HotDailyCard), incrementIf(selectedDailyCard.CardType == study.HotDailyCard, studyProgress.HotCardsleft),
		underlineIf(selectedDailyCard.CardType == study.NewDailyCard), incrementIf(selectedDailyCard.CardType == study.NewDailyCard, studyProgress.NewCardsLeft),
		underlineIf(selectedDailyCard.CardType == study.ScheduledDailyCard), incrementIf(selectedDailyCard.CardType == study.ScheduledDailyCard, studyProgress.ScheduledCardsLeft))

	fmt.Print(sgr.MustParseln(str))
}

func (i *CommandLineInteractor) AskQuestion(crd *card.Card, states []*study.CardState) (*study.CardState, error) {
	i.showQuestionHeader(crd)
	i.showQuestion(crd)
	err := i.waitForUserToComeUpWithAnswer(crd)
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

func (i *CommandLineInteractor) waitForUserToComeUpWithAnswer(crd *card.Card) error {
	fmt.Println(sgr.MustParseln("[fg-green][underline]C[underlineOff]opy path [fg-yellow][underline]S[underlineOff]how answer[reset]"))

	for {
		ch, key, err := keyboard.GetSingleKey()
		if err != nil {
			return err
		}

		if ch == 'c' || ch == 'C' {
			i.copyFullCardPathToClipboard(crd)
		} else if key == 13 || ch == 's' || ch == 'S' {
			return nil
		} else if key == 3 || key == 27 {
			// ctrl + c or escape
			return io.EOF
		}
	}
}

func (i *CommandLineInteractor) copyFullCardPathToClipboard(crd *card.Card) {
	clipboard.WriteAll(fmt.Sprintf("%s/%s", crd.Section(), crd.Entry()))
	fmt.Println("Card path has been copied to clipboard")
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

		var interval string

		if state.Interval.Minutes() < 60 {
			interval = strconv.Itoa(int(state.Interval.Minutes())) + "m"
		} else if state.Interval.Hours() < 24 {
			interval = strconv.Itoa(int(state.Interval.Hours())) + "d"
		} else {
			interval = strconv.Itoa(int(state.Interval.Hours()/24)) + "d"
		}

		fmt.Printf(sgr.MustParse(gradeColor+"[underline]%s[underlineOff]%s (%s) "), string(state.Name[0]), state.Name[1:], interval)
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
