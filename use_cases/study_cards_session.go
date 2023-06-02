package use_cases

import (
	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/core/study"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type StudyCardsSession struct {
	timeSource     datetime.TimeSource
	cardRepo       card.CardRepository
	userInteractor study.UserInteractor
}

func NewStudyCardsSession(timeSource datetime.TimeSource, cardRepo card.CardRepository, userInteractor study.UserInteractor) *StudyCardsSession {
	return &StudyCardsSession{
		timeSource:     timeSource,
		cardRepo:       cardRepo,
		userInteractor: userInteractor,
	}
}

func (s *StudyCardsSession) Run(cardEmitter pipeline.Emitter[card.BriefCard], config *study.DailyCardsConfig) error {
	session := study.NewSession(s.timeSource, s.cardRepo, s.userInteractor)

	session.Run(config, cardEmitter)

	return nil
}
