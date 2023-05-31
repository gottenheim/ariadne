package study

import (
	"context"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type newCardsCollector struct {
	timeSource datetime.TimeSource
	cardRepo   card.CardRepository
	config     *DailyCardsConfig
	newCards   []*card.Card
}

func CollectNewCards(timeSource datetime.TimeSource, cardRepo card.CardRepository, config *DailyCardsConfig) *newCardsCollector {
	return &newCardsCollector{
		timeSource: timeSource,
		cardRepo:   cardRepo,
		config:     config,
	}
}

func (f *newCardsCollector) Run(ctx context.Context, input <-chan card.BriefCard, output chan<- card.BriefCard) error {
	var newCards, cardsLearnedToday []*card.Card

	for {
		briefCard, ok := <-input
		if !ok {
			break
		}

		isNewOrLearnedToday, err := f.isCardNewOrLearnedToday(briefCard)
		if err != nil {
			return err
		}

		if !isNewOrLearnedToday {
			pipeline.WriteToChannel[card.BriefCard](ctx, output, briefCard)
			continue
		}

		crd, err := f.cardRepo.Get(briefCard.Section, briefCard.Entry)

		if err != nil {
			return err
		}

		isNewCard, err := card.IsNewCardActivities(briefCard.Activities)

		if err != nil {
			return err
		}

		if isNewCard {
			newCards = append(newCards, crd)
		} else {
			cardsLearnedToday = append(cardsLearnedToday, crd)
		}
	}

	newCardsRemaining := f.config.NewCardsCount - len(cardsLearnedToday)
	if newCardsRemaining < 0 {
		newCardsRemaining = 0
	}

	if len(newCards) < newCardsRemaining {
		newCardsRemaining = len(newCards)
	}

	f.newCards = newCards[0:newCardsRemaining]

	return nil
}

func (f *newCardsCollector) GetNewCards() []*card.Card {
	return f.newCards
}

func (f *newCardsCollector) isCardNewOrLearnedToday(briefCard card.BriefCard) (bool, error) {
	isNewCard, err := card.IsNewCardActivities(briefCard.Activities)

	if err != nil {
		return false, err
	}

	if isNewCard {
		return true, nil
	}

	isCardLearnedToday, err := card.IsCardLearnedToday(f.timeSource, briefCard.Activities)

	if err != nil {
		return false, err
	}

	return isCardLearnedToday, nil
}
