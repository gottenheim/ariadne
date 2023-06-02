package interactor

import (
	"fmt"

	sgr "github.com/foize/go.sgr"
	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
)

type CommandLineInteractor struct {
}

func NewCommandLineInteractor() study.UserInteractor {
	return &CommandLineInteractor{}
}

func (i *CommandLineInteractor) ShowDiscoveredDailyCards(dailyCards *study.DailyCards) {
	fmt.Printf(sgr.MustParseln("[fg-red]%d new cards[reset], [fg-green]%d new cards[reset], [fg-blue]%d cards to remind[reset]"),
		len(dailyCards.HotCardsToRevise), len(dailyCards.NewCards), len(dailyCards.ScheduledCards))
}

func (i *CommandLineInteractor) AskQuestion(crd *card.Card, states []*study.CardState) (*study.CardState, error) {
	fmt.Println("----------- Question -------------")
	fmt.Printf(sgr.MustParseln("Directory: [fg-white]%s/%s\n"), crd.Section(), crd.Entry())

	fmt.Println(sgr.MustParseln("[fg-white][underline]S[underlineOff]how answer[reset]"))

	return nil, nil
}
