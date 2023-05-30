package study

import (
	"context"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type newCardCondition struct {
	timeSource datetime.TimeSource
	cardRepo   card.CardRepository
}

func NewCardCondition(timeSource datetime.TimeSource, cardRepo card.CardRepository) pipeline.Condition[card.BriefCard, *card.Card, card.BriefCard] {
	return &newCardCondition{
		timeSource: timeSource,
		cardRepo:   cardRepo,
	}
}

func (f *newCardCondition) Run(ctx context.Context, input <-chan card.BriefCard, positiveDecision chan<- *card.Card, negativeDecision chan<- card.BriefCard) error {
	for {
		briefCard, ok := <-input

		if !ok {
			break
		}

		isNewCard, err := card.IsNewCardActivities(briefCard.Activities)

		if err != nil {
			return err
		}

		isCardLearnedToday, err := card.IsCardLearnedToday(f.timeSource, briefCard.Activities)

		if err != nil {
			return err
		}

		if isNewCard || isCardLearnedToday {
			card, err := f.cardRepo.Get(briefCard.Key)
			if err != nil {
				return err
			}

			positiveDecision <- card
		} else {
			negativeDecision <- briefCard
		}
	}

	return nil
}
