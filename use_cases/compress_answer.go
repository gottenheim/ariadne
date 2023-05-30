package use_cases

import "github.com/gottenheim/ariadne/core/card"

type CompressAnswerAction struct {
}

func (a *CompressAnswerAction) Run(cardRepo card.CardRepository, cardKey card.Key) error {
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
