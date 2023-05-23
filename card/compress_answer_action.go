package card

import (
	"github.com/spf13/afero"
)

type CompressAnswerAction struct {
}

func (a *CompressAnswerAction) Run(fs afero.Fs, baseDirPath string, cardDirPath string) error {
	cardRepo := NewFileCardRepository(fs, baseDirPath)

	card, err := cardRepo.Get(cardDirPath)
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
