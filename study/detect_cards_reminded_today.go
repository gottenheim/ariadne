package study

import (
	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
)

type detectCardsRemindedToday struct {
	timeSource datetime.TimeSource
}

func (r *detectCardsRemindedToday) Run(filter cardFilterInterface, cardKey string, cardActivity card.CardActivity, next actionFunc) error {
	if !filter.needCardsToRemind() {
		return next(filter, cardKey, cardActivity)
	}

	isRemindedToday, err := card.IsCardRemindedToday(r.timeSource, cardActivity)

	if err != nil {
		return err
	}

	if !isRemindedToday {
		return next(filter, cardKey, cardActivity)
	}

	filter.decrementCardsToRemindCounter()

	return nil
}
