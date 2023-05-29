package card

import (
	"context"

	"github.com/gottenheim/ariadne/details/datetime"
	"github.com/gottenheim/ariadne/details/pipeline"
)

type scheduledCardCondition struct {
	timeSource datetime.TimeSource
	cardRepo   CardRepository
}

func ScheduledCardCondition(timeSource datetime.TimeSource, cardRepo CardRepository) pipeline.Condition[BriefCard, *Card, BriefCard] {
	return &scheduledCardCondition{
		timeSource: timeSource,
		cardRepo:   cardRepo,
	}
}

func (f *scheduledCardCondition) Run(ctx context.Context, input <-chan BriefCard, positiveDecision chan<- *Card, negativeDecision chan<- BriefCard) error {
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
