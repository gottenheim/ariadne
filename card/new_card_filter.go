package card

import "github.com/gottenheim/ariadne/details/pipeline"

type KeyWithActivities struct {
	Key        Key
	Activities CardActivity
}

type newCardFilter struct {
	cardRepo CardRepository
}

func NewCardFilter(cardRepo CardRepository) pipeline.Filter[*KeyWithActivities, *Card] {
	return &newCardFilter{
		cardRepo: cardRepo,
	}
}

func (f *newCardFilter) Run(input <-chan *KeyWithActivities, output chan<- *Card) error {
	defer func() {
		close(output)
	}()

	for {
		keyWithActivities, ok := <-input

		if !ok {
			break
		}

		isNewCard, err := IsNewCard(keyWithActivities.Activities)

		if err != nil {
			return err
		}

		if isNewCard {
			card, err := f.cardRepo.Get(keyWithActivities.Key)
			if err != nil {
				return err
			}

			output <- card
		}
	}

	return nil
}
