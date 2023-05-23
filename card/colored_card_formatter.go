package card

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

type ColoredCardFormatter struct {
	writer io.Writer
}

func NewColoredCardFormatter(writer io.Writer) *ColoredCardFormatter {
	return &ColoredCardFormatter{
		writer: writer,
	}
}

func (f *ColoredCardFormatter) FormatCard(card *Card) error {
	for _, artifact := range card.artifacts {
		if artifact.name == AnswerArtifactName {
			continue
		}

		_, err := f.writer.Write([]byte(color.GreenString("----- %s -----\n", artifact.name)))
		if err != nil {
			return err
		}

		_, err = f.writer.Write([]byte(fmt.Sprintf("%s\n\n", artifact.content)))
		if err != nil {
			return err
		}
	}

	return nil
}
