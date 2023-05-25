package card

type NewCardAction struct {
}

func (a *NewCardAction) Run(templateRepo CardTemplateRepository, cardRepo CardRepository, cardSection []string) error {
	cardTemplate, err := templateRepo.GetTemplate()
	if err != nil {
		return err
	}

	card := NewCard(cardSection, 0, cardTemplate.Artifacts())

	return cardRepo.Save(card)
}
