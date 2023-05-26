package card

type NewCardAction struct {
}

func (a *NewCardAction) Run(templateRepo CardTemplateRepository, cardRepo CardRepository) error {
	cardTemplate, err := templateRepo.GetTemplate()
	if err != nil {
		return err
	}

	card := NewCard(0, cardTemplate.Artifacts())

	return cardRepo.Save(card)
}
