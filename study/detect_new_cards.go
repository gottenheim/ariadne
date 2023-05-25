package study

import (
	"github.com/gottenheim/ariadne/card"
)

type detectNewCards struct {
	cardRepo card.CardRepository
}

func (r *detectNewCards) Run(filter cardFilterInterface, cardKey string, cardActivity card.CardActivity, next actionFunc) error {
	if !filter.needNewCards() {
		return next(filter, cardKey, cardActivity)
	}

	isNewCard, err := card.IsNewCard(cardActivity)

	if err != nil {
		return err
	}

	if !isNewCard {
		return next(filter, cardKey, cardActivity)
	}

	card, err := r.cardRepo.Get(cardKey)
	if err != nil {
		return err
	}

	filter.addCard(card)
	filter.decrementNewCardsCounter()

	return nil
}
