package use_cases

import (
	"io"

	"github.com/gottenheim/ariadne/core/card"
)

type ShowAnswer struct {
}

func (a *ShowAnswer) Run(cardRepo card.CardRepository, output io.Writer, cardKey card.Key) error {
	c, err := cardRepo.Get(cardKey)

	if err != nil {
		return err
	}

	formatter := card.NewColoredCardFormatter(output)

	formatter.FormatCard(c)

	return nil
}
