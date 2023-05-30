package use_cases

import "github.com/gottenheim/ariadne/core/card"

type ExtractCard struct {
}

func (a *ExtractCard) Run(cardRepo card.CardRepository, cardKey card.Key) error {
	card, err := cardRepo.Get(cardKey)
	if err != nil {
		return err
	}

	err = card.ExtractAnswer()
	if err != nil {
		return err
	}

	return cardRepo.Save(card)
}
