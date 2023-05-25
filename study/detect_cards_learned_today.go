package study

import (
	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
)

type detectCardsLearnedToday struct {
	timeSource datetime.TimeSource
}

func (r *detectCardsLearnedToday) Run(filter cardFilterInterface, cardKey string, cardActivity card.CardActivity, next actionFunc) error {
	if !filter.needNewCards() {
		return next(filter, cardKey, cardActivity)
	}

	isLearnedToday, err := card.IsCardLearnedToday(r.timeSource, cardActivity)

	if err != nil {
		return err
	}

	if !isLearnedToday {
		return next(filter, cardKey, cardActivity)
	}

	filter.decrementNewCardsCounter()

	return nil
}
