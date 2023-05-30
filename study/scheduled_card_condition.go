package study

import (
	"context"

	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
	"github.com/gottenheim/ariadne/details/pipeline"
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

		isCardScheduledToRemindToday, err := IsCardScheduledToRemindToday(f.timeSource, briefCard.Activities)

		if err != nil {
			return err
		}

		isCardRemindedToday, err := IsCardRemindedToday(f.timeSource, briefCard.Activities)

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
