package card

import (
	"github.com/spf13/afero"
)

type ExtractCardAction struct {
}

func (a *ExtractCardAction) Run(fs afero.Fs, baseDirPath string, cardDirPath string) error {
	cardRepo := NewFileCardRepository(fs, baseDirPath)

	card, err := cardRepo.Get(cardDirPath)
	if err != nil {
		return err
	}

	err = card.ExtractAnswer()
	if err != nil {
		return err
	}

	return cardRepo.Save(card)
}
