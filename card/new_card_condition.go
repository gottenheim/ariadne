package card

import "github.com/gottenheim/ariadne/details/pipeline"

type newCardCondition struct {
	cardRepo CardRepository
}

func NewCardCondition(cardRepo CardRepository) pipeline.Condition[BriefCard, *Card, BriefCard] {
	return &newCardCondition{
		cardRepo: cardRepo,
	}
}

func (f *newCardCondition) Run(input <-chan BriefCard, positiveDecision chan<- *Card, negativeDecision chan<- BriefCard) error {
	defer func() {
		close(positiveDecision)
		close(negativeDecision)
	}()

	for {
		briefCard, ok := <-input

		if !ok {
			break
		}

		isNewCard, err := IsNewCardActivities(briefCard.Activities)

		if err != nil {
			return err
		}

		if isNewCard {
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
