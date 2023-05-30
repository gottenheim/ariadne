package study

import (
	"context"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/datetime"
	"github.com/gottenheim/ariadne/libraries/pipeline"
)

type scheduledCardCondition struct {
	timeSource datetime.TimeSource
	cardRepo   card.CardRepository
}

func ScheduledCardCondition(timeSource datetime.TimeSource, cardRepo card.CardRepository) pipeline.Condition[card.BriefCard, *card.Card, card.BriefCard] {
	return &scheduledCardCondition{
		timeSource: timeSource,
		cardRepo:   cardRepo,
	}
}

func (f *scheduledCardCondition) Run(ctx context.Context, input <-chan card.BriefCard, positiveDecision chan<- *card.Card, negativeDecision chan<- card.BriefCard) error {
	for {
		briefCard, ok := <-input

		if !ok {
			break
		}

		isCardScheduledToRemindToday, err := card.IsCardScheduledToRemindToday(f.timeSource, briefCard.Activities)

		if err != nil {
			return err
		}

		isCardRemindedToday, err := card.IsCardRemindedToday(f.timeSource, briefCard.Activities)

		if err != nil {
			return err
		}

		if isCardScheduledToRemindToday || isCardRemindedToday {
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
