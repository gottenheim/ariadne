package card

import (
	"os"
	"path"

	"github.com/spf13/afero"
)

type FileTemplateRepository struct {
	fs          afero.Fs
	templateDir string
}

func NewFileTemplateRepository(fs afero.Fs, templateDir string) CardTemplateRepository {
	return &FileTemplateRepository{
		fs:          fs,
		templateDir: templateDir,
	}
}

func (r *FileTemplateRepository) GetTemplate() (*CardTemplate, error) {
	var artifacts []CardArtifact

	err := afero.Walk(r.fs, r.templateDir, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileContents, err := afero.ReadFile(r.fs, filePath)
			if err != nil {
				return err
			}
			fileName := path.Base(filePath)
			artifacts = append(artifacts, NewCardArtifact(fileName, fileContents))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return NewCardTemplate(artifacts), nil
}
