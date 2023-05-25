package card_repo

import (
	"os"
	"path"

	"github.com/gottenheim/ariadne/card"
	"github.com/spf13/afero"
)

type FileTemplateRepository struct {
	fs          afero.Fs
	templateDir string
}

func NewFileTemplateRepository(fs afero.Fs, templateDir string) card.CardTemplateRepository {
	return &FileTemplateRepository{
		fs:          fs,
		templateDir: templateDir,
	}
}

func (r *FileTemplateRepository) GetTemplate() (*card.CardTemplate, error) {
	var artifacts []card.CardArtifact

	err := afero.Walk(r.fs, r.templateDir, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileContents, err := afero.ReadFile(r.fs, filePath)
			if err != nil {
				return err
			}
			fileName := path.Base(filePath)
			artifacts = append(artifacts, card.NewCardArtifact(fileName, fileContents))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return card.NewCardTemplate(artifacts), nil
}
