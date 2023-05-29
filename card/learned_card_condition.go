package card

import (
	"github.com/gottenheim/ariadne/details/datetime"
	"github.com/gottenheim/ariadne/details/pipeline"
)

type learnedCardCondition struct {
	timeSource datetime.TimeSource
	cardRepo   CardRepository
}

func LearnedCardCondition(timeSource datetime.TimeSource, cardRepo CardRepository) pipeline.Condition[BriefCard, *Card, BriefCard] {
	return &learnedCardCondition{
		timeSource: timeSource,
		cardRepo:   cardRepo,
	}
}

func (f *learnedCardCondition) Run(input <-chan BriefCard, positiveDecision chan<- *Card, negativeDecision chan<- BriefCard) error {
	defer func() {
		close(positiveDecision)
		close(negativeDecision)
	}()

	for {
		briefCard, ok := <-input

		if !ok {
			break
		}

		isCardScheduledToRemind, err := IsCardLearnedToday(f.timeSource, briefCard.Activities)

		if err != nil {
			return err
		}

		if isCardScheduledToRemind {
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
