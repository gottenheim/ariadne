package interactor

import (
	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
)

type FakeUserInteractor struct {
}

func NewFakeUserInteractor() study.UserInteractor {
	return &FakeUserInteractor{}
}

func (i *FakeUserInteractor) ShowDiscoveredDailyCards(dailyCards *study.DailyCards) {
}

func (i *FakeUserInteractor) AskQuestion(crd *card.Card, states []*study.CardState) (*study.CardState, error) {
	return nil, nil
}
