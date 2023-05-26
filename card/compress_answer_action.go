package card

type CompressAnswerAction struct {
}

func (a *CompressAnswerAction) Run(cardRepo CardRepository, cardKey Key) error {
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
