package card

type ExtractCardAction struct {
}

func (a *ExtractCardAction) Run(cardRepo CardRepository, cardKey string) error {
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
