package study

import (
	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
)

type cardFilterInterface interface {
	needNewCards() bool
	needCardsToRemind() bool
	addCard(card *card.Card)
	decrementNewCardsCounter()
	decrementCardsToRemindCounter()
}

type cardFilter struct {
	repo          card.CardRepository
	newCards      int
	cardsToRemind int
	actions       []cardFilterAction
	cards         []*card.Card
}

func NewCardFilter(cardRepo card.CardRepository, timeSource datetime.TimeSource, newCards int, cardsToRemind int) *cardFilter {
	return &cardFilter{
		actions: []cardFilterAction{
			&detectNewCards{
				cardRepo: cardRepo,
			},
			&detectCardsLearnedToday{
				timeSource: timeSource,
			},
			&detectCardsToRemind{
				timeSource: timeSource,
				cardRepo:   cardRepo,
			},
			&detectCardsRemindedToday{
				timeSource: timeSource,
			},
		},
	}
}

func (f *cardFilter) ProcessCard(cardKey string, cardActivity card.CardActivity) error {
	nextAction := func(filter cardFilterInterface, cardKey string, cardActivity card.CardActivity) error {
		return nil
	}
	for i := len(f.actions) - 1; i > 0; i-- {
		action := f.actions[i]
		runAction := func(filter cardFilterInterface, cardKey string, cardActivity card.CardActivity) error {
			return action.Run(filter, cardKey, cardActivity, nextAction)
		}
		nextAction = runAction
	}

	return f.actions[0].Run(f, cardKey, cardActivity, nextAction)
}

func (f *cardFilter) IsDone() bool {
	return f.cardsToRemind == 0 && f.newCards == 0
}

func (f *cardFilter) needNewCards() bool {
	return f.newCards > 0
}

func (f *cardFilter) needCardsToRemind() bool {
	return f.cardsToRemind > 0
}

func (f *cardFilter) addCard(card *card.Card) {
	f.cards = append(f.cards, card)
}

func (f *cardFilter) decrementNewCardsCounter() {
	f.newCards--
}

func (f *cardFilter) decrementCardsToRemindCounter() {
	f.newCards--
}
