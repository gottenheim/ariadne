package interactor

import (
	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
)

type ChooseStateFunc func(*card.Card, []*study.CardState) (*study.CardState, error)

type FakeUserInteractor struct {
	chooseState ChooseStateFunc
}

func NewFakeUserInteractor(chooseState ChooseStateFunc) *FakeUserInteractor {
	return &FakeUserInteractor{
		chooseState: chooseState,
	}
}

func (i *FakeUserInteractor) ShowDiscoveredDailyCards(dailyCards *study.DailyCards) {
}

func (i *FakeUserInteractor) AskQuestion(crd *card.Card, states []*study.CardState) (*study.CardState, error) {
	return i.chooseState(crd, states)
}
