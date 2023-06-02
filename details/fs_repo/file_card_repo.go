package fs_repo

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/fs"
	"github.com/spf13/afero"
)

type fileCardRepository struct {
	fs afero.Fs
}

func NewFileCardRepository(fs afero.Fs) *fileCardRepository {
	return &fileCardRepository{
		fs: fs,
	}
}

func (r *fileCardRepository) Get(section string, entry string) (*card.Card, error) {
	cardPath := r.GetCardPath(section, entry)

	artifacts, err := r.readCardArtifacts(cardPath)
	if err != nil {
		return nil, err
	}

	if len(artifacts) == 0 {
		return nil, errors.New("Path %s doesn't contain card artifacts")
	}

	activities, err := r.ReadCardActivities(cardPath)
	if err != nil {
		return nil, err
	}

	return card.RestoreExisting(section, entry, artifacts, activities), nil
}

func (r *fileCardRepository) Save(card *card.Card) error {
	err := r.generateKeyIfNeeded(card)
	if err != nil {
		return err
	}

	err = r.clearCardDirectory(card)
	if err != nil {
		return err
	}

	err = r.saveArtifactFiles(card)
	if err != nil {
		return err
	}

	return r.SaveCardActivities(card.Activities(), r.GetCardPath(card.Section(), card.Entry()))
}

func (r *fileCardRepository) generateKeyIfNeeded(card *card.Card) error {
	if len(card.Entry()) == 0 {
		freeEntry, err := r.getNextFreeSectionEntry(card.Section())

		if err != nil {
			return err
		}

		card.SetEntry(freeEntry)
	}

	return nil
}

func (r *fileCardRepository) getNextFreeSectionEntry(section string) (string, error) {
	maxCardKey := 0

	err := afero.Walk(r.fs, section, func(filePath string, info os.FileInfo, err error) error {
		isDir, _ := afero.IsDir(r.fs, filePath)

		if isDir {
			cardDir := path.Base(filePath)
			cardKey, err := strconv.Atoi(cardDir)
			if err != nil {
				return nil
			}

			if cardKey > maxCardKey {
				maxCardKey = cardKey
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return strconv.Itoa(maxCardKey + 1), nil
}

func (r *fileCardRepository) GetCardPath(section string, entry string) string {
	return filepath.Join(section, entry)
}

func (r *fileCardRepository) GetCardPathSection(cardPath string) string {
	return filepath.Dir(cardPath)
}

func (r *fileCardRepository) GetCardPathEntry(cardPath string) string {
	return filepath.Base(cardPath)
}

func (r *fileCardRepository) clearCardDirectory(card *card.Card) error {
	cardPath := r.GetCardPath(card.Section(), card.Entry())

	cardDirExists, err := afero.Exists(r.fs, cardPath)

	if err != nil {
		return err
	}

	if cardDirExists {
		fs.RemoveAllDirectoryFiles(r.fs, cardPath)
	}

	return nil
}

func (r *fileCardRepository) saveArtifactFiles(card *card.Card) error {
	cardPath := r.GetCardPath(card.Section(), card.Entry())

	err := r.fs.MkdirAll(cardPath, os.ModePerm)

	if err != nil {
		return err
	}

	for _, artifact := range card.Artifacts() {
		filePath := filepath.Join(cardPath, artifact.Name())
		err = afero.WriteFile(r.fs, filePath, artifact.Content(), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *fileCardRepository) readCardArtifacts(cardPath string) ([]card.CardArtifact, error) {
	var artifacts []card.CardArtifact

	err := afero.Walk(r.fs, cardPath, func(filePath string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && !r.isServiceFile(filePath) {
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

	return artifacts, nil
}

func (r *fileCardRepository) isServiceFile(fileName string) bool {
	return fileName[0] == '.'
}
