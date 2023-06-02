package use_cases

import (
	"github.com/gottenheim/ariadne/core/card"
)

type ShowAnswer struct {
}

func (a *ShowAnswer) Run(cardRepo card.CardRepository, userInteractor card.UserInteractor, section string, entry string) error {
	crd, err := cardRepo.Get(section, entry)

	if err != nil {
		return err
	}

	userInteractor.ShowAnswer(crd)

	return nil
}
