package card

import (
	"io"

	"github.com/spf13/afero"
)

type ShowAnswerAction struct {
}

func (a *ShowAnswerAction) Run(fs afero.Fs, writer io.Writer, baseDirPath string, cardDirPath string) error {
	cardRepo := NewFileCardRepository(fs, baseDirPath)

	card, err := cardRepo.Get(cardDirPath)

	if err != nil {
		return err
	}

	formatter := NewColoredCardFormatter(writer)

	formatter.FormatCard(card)

	return nil
}
