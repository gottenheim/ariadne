package study

import (
	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
)

type detectCardsToRemind struct {
	timeSource datetime.TimeSource
	cardRepo   card.CardRepository
}

func (r *detectCardsToRemind) Run(filter cardFilterInterface, cardKey string, cardActivity card.CardActivity, next actionFunc) error {
	if !filter.needCardsToRemind() {
		return next(filter, cardKey, cardActivity)
	}

	isScheduledToRemind, err := card.IsCardScheduledToRemindToday(r.timeSource, cardActivity)

	if err != nil {
		return err
	}

	if !isScheduledToRemind {
		return next(filter, cardKey, cardActivity)
	}

	card, err := r.cardRepo.Get(cardKey)
	if err != nil {
		return err
	}

	filter.decrementCardsToRemindCounter()
	filter.addCard(card)

	return nil
}
