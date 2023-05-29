package card

import (
	"context"

	"github.com/gottenheim/ariadne/details/datetime"
	"github.com/gottenheim/ariadne/details/pipeline"
)

type remindedCardCondition struct {
	timeSource datetime.TimeSource
	cardRepo   CardRepository
}

func RemindedCardCondition(timeSource datetime.TimeSource, cardRepo CardRepository) pipeline.Condition[BriefCard, *Card, BriefCard] {
	return &remindedCardCondition{
		timeSource: timeSource,
		cardRepo:   cardRepo,
	}
}

func (f *remindedCardCondition) Run(ctx context.Context, input <-chan BriefCard, positiveDecision chan<- *Card, negativeDecision chan<- BriefCard) error {
	for {
		briefCard, ok := <-input

		if !ok {
			break
		}

		isCardRemindedToday, err := IsCardRemindedToday(f.timeSource, briefCard.Activities)

		if err != nil {
			return err
		}

		if isCardRemindedToday {
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
