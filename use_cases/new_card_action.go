package use_cases

import "github.com/gottenheim/ariadne/core/card"

type NewCard struct {
}

func (a *NewCard) Run(templateRepo card.CardTemplateRepository, cardRepo card.CardRepository, section string) error {
	cardTemplate, err := templateRepo.GetTemplate()
	if err != nil {
		return err
	}

	card := card.CreateNew(section, cardTemplate.Artifacts())

	return cardRepo.Save(card)
}
