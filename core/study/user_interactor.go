package study

import (
	"github.com/gottenheim/ariadne/core/card"
)

type UserInteractor interface {
	ShowDiscoveredDailyCards(dailyCards *DailyCards)
	AskQuestion(crd *card.Card, states []*CardState) (*CardState, error)
}
