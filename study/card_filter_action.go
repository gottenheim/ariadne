package study

import "github.com/gottenheim/ariadne/card"

type actionFunc func(cardFilterInterface, string, card.CardActivity) error

type cardFilterAction interface {
	Run(filter cardFilterInterface, cardKey string, cardActivity card.CardActivity, next actionFunc) error
}
