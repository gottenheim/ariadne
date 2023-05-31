package use_cases

import (
	"io"

	"github.com/gottenheim/ariadne/core/card"
)

type ShowAnswer struct {
}

func (a *ShowAnswer) Run(cardRepo card.CardRepository, output io.Writer, section string, entry string) error {
	c, err := cardRepo.Get(section, entry)

	if err != nil {
		return err
	}

	formatter := card.NewColoredCardFormatter(output)

	formatter.FormatCard(c)

	return nil
}
