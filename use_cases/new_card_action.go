package use_cases

import "github.com/gottenheim/ariadne/core/card"

type NewCardAction struct {
}

func (a *NewCardAction) Run(templateRepo card.CardTemplateRepository, cardRepo card.CardRepository) error {
	cardTemplate, err := templateRepo.GetTemplate()
	if err != nil {
		return err
	}

	card := card.NewCard(0, cardTemplate.Artifacts())

	return cardRepo.Save(card)
}
