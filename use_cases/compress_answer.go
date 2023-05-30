package use_cases

import "github.com/gottenheim/ariadne/core/card"

type CompressAnswer struct {
}

func (a *CompressAnswer) Run(cardRepo card.CardRepository, cardKey card.Key) error {
	card, err := cardRepo.Get(cardKey)
	if err != nil {
		return err
	}

	err = card.CompressAnswer()
	if err != nil {
		return err
	}

	err = cardRepo.Save(card)
	if err != nil {
		return err
	}

	return nil
}
