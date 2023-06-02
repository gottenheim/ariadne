package study

import (
	"github.com/go-errors/errors"
	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type ChooseStateFunc func(*card.Card, []*CardState) (*CardState, error)

type Session struct {
	timeSource datetime.TimeSource
	cardRepo   card.CardRepository
}

func NewSession(timeSource datetime.TimeSource, cardRepo card.CardRepository) *Session {
	return &Session{
		timeSource: timeSource,
		cardRepo:   cardRepo,
	}
}

func (s *Session) Run(dailyCardsConfig *DailyCardsConfig, cardEmitter pipeline.Emitter[card.BriefCard], chooseState ChooseStateFunc) error {
	dailyCards, err := s.collectDailyCards(dailyCardsConfig, cardEmitter)

	if err != nil {
		return err
	}

	return s.studyDailyCards(dailyCards, chooseState)
}

func (s *Session) collectDailyCards(dailyCardsConfig *DailyCardsConfig, cardEmitter pipeline.Emitter[card.BriefCard]) (*DailyCards, error) {
	cardsCollector := NewDailyCardsCollector(s.timeSource, s.cardRepo, cardEmitter)
	cardsCollector.SetConfig(dailyCardsConfig)
	return cardsCollector.Collect()
}

func (s *Session) studyDailyCards(dailyCards *DailyCards, chooseState ChooseStateFunc) error {
	it := NewDailyCardsIterator(s.timeSource, dailyCards)

	for {
		crd, err := it.Next()

		if err != nil {
			return err
		}

		if crd == nil {
			return nil
		}

		err = s.moveCardToNextState(crd, chooseState)

		if err != nil {
			return err
		}

		err = s.reAttachCardToSessionIfItRemainsScheduledToToday(crd, it)

		if err != nil {
			return err
		}
	}
}

func (s *Session) moveCardToNextState(crd *card.Card, chooseState ChooseStateFunc) error {
	cardWorkflow := NewCardWorkflow(s.timeSource, crd)

	nextStates, err := cardWorkflow.GetNextStates()

	if err != nil {
		return err
	}

	chosenState, err := chooseState(crd, nextStates)

	if err != nil {
		return err
	}

	if chosenState == nil {
		return errors.New("No state was chosen")
	}

	err = cardWorkflow.TransitTo(chosenState)

	if err != nil {
		return err
	}

	s.cardRepo.Save(crd)

	return nil
}

func (s *Session) reAttachCardToSessionIfItRemainsScheduledToToday(crd *card.Card, it *DailyCardsIterator) error {
	isScheduledToRemindToday, err := card.IsCardScheduledToRemindToday(s.timeSource, crd.Activities())

	if err != nil {
		return err
	}

	if isScheduledToRemindToday {
		it.AddHotCardToRevise(crd)
	}

	return nil
}
