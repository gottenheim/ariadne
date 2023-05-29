package card

import (
	"context"

	"github.com/gottenheim/ariadne/details/datetime"
	"github.com/gottenheim/ariadne/details/pipeline"
)

type newCardCondition struct {
	timeSource datetime.TimeSource
	cardRepo   CardRepository
}

func NewCardCondition(timeSource datetime.TimeSource, cardRepo CardRepository) pipeline.Condition[BriefCard, *Card, BriefCard] {
	return &newCardCondition{
		timeSource: timeSource,
		cardRepo:   cardRepo,
	}
}

func (f *newCardCondition) Run(ctx context.Context, input <-chan BriefCard, positiveDecision chan<- *Card, negativeDecision chan<- BriefCard) error {
	for {
		briefCard, ok := <-input

		if !ok {
			break
		}

		isNewCard, err := IsNewCardActivities(briefCard.Activities)

		if err != nil {
			return err
		}

		isCardLearnedToday, err := IsCardLearnedToday(f.timeSource, briefCard.Activities)

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
