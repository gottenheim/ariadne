package card

import "github.com/gottenheim/ariadne/details/pipeline"

type KeyWithActivities struct {
	Key        Key
	Activities CardActivity
}

type newCardFilter struct {
	events   pipeline.FilterEvents
	cardRepo CardRepository
}

func NewCardFilter(events pipeline.FilterEvents, cardRepo CardRepository) pipeline.Filter[*KeyWithActivities, *Card] {
	events.OnStart()
	return &newCardFilter{
		events:   events,
		cardRepo: cardRepo,
	}
}

func (f *newCardFilter) Run(input <-chan *KeyWithActivities, output chan<- *Card) {
	defer func() {
		close(output)
		f.events.OnFinish()
	}()

	for {
		keyWithActivities, ok := <-input

		if !ok {
			break
		}

		isNewCard, err := IsNewCard(keyWithActivities.Activities)

		if err != nil {
			f.events.OnError(err)
			break
		}

		if isNewCard {
			card, err := f.cardRepo.Get(keyWithActivities.Key)
			if err != nil {
				f.events.OnError(err)
				break
			}

			output <- card
		}
	}
}
