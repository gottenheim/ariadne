package card

import (
	"io"
)

type ShowAnswerAction struct {
}

func (a *ShowAnswerAction) Run(cardRepo CardRepository, output io.Writer, cardKey string) error {
	card, err := cardRepo.Get(cardKey)

	if err != nil {
		return err
	}

	formatter := NewColoredCardFormatter(output)

	formatter.FormatCard(card)

	return nil
}
