package card

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type NewCardAction struct {
}

func (a *NewCardAction) Run(fs afero.Fs, baseDirPath string, cardsDirPath string, templateDirPath string) error {
	templateRepo := NewFileTemplateRepository(fs, templateDirPath)
	cardTemplate, err := templateRepo.GetTemplate()
	if err != nil {
		return err
	}

	cardSections := a.getCardSections(cardsDirPath)

	card := NewCard(cardSections, 0, cardTemplate.Artifacts())

	cardRepo := NewFileCardRepository(fs, baseDirPath)

	return cardRepo.Save(card)
}

func (a *NewCardAction) getCardSections(cardsDirPath string) []string {
	return filepath.SplitList(cardsDirPath)
}
