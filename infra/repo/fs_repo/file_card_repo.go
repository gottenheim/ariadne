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
	fs       afero.Fs
	cardsDir string
}

func NewFileCardRepository(fs afero.Fs, cardsDir string) card.CardRepository {
	return &fileCardRepository{
		fs:       fs,
		cardsDir: cardsDir,
	}
}

func newFileCardRepository(fs afero.Fs, cardsDir string) *fileCardRepository {
	return &fileCardRepository{
		fs:       fs,
		cardsDir: cardsDir,
	}
}

func (r *fileCardRepository) Get(cardKey card.Key) (*card.Card, error) {
	card, err := r.createCard(cardKey)
	if err != nil {
		return nil, err
	}

	cardPath := r.getCardPath(card)

	artifacts, err := r.readCardArtifacts(cardPath)
	if err != nil {
		return nil, err
	}

	if len(artifacts) == 0 {
		return nil, errors.New("Path %s doesn't contain card artifacts")
	}

	card.SetArtifacts(artifacts)

	activities, err := r.ReadCardActivities(cardPath)
	if err != nil {
		return nil, err
	}

	card.SetActivities(activities)

	return card, nil
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

	return r.SaveCardActivities(card.Activities(), r.getCardPath(card))
}

func (r *fileCardRepository) generateKeyIfNeeded(card *card.Card) error {
	if card.Key() == 0 {
		cardKey, err := r.getNextFreeCardKey()

		if err != nil {
			return err
		}

		card.SetKey(cardKey)
	}

	return nil
}

func (r *fileCardRepository) getNextFreeCardKey() (card.Key, error) {
	maxCardKey := 0

	err := afero.Walk(r.fs, r.cardsDir, func(filePath string, info os.FileInfo, err error) error {
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
		return 0, err
	}

	return card.Key(maxCardKey + 1), nil
}

func (r *fileCardRepository) getCardPath(card *card.Card) string {
	return filepath.Join(r.cardsDir, strconv.Itoa(int(card.Key())))
}

func (r *fileCardRepository) clearCardDirectory(card *card.Card) error {
	cardPath := r.getCardPath(card)

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
	cardPath := r.getCardPath(card)

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

func (r *fileCardRepository) createCard(cardKey card.Key) (*card.Card, error) {
	return card.NewCard(cardKey, []card.CardArtifact{}), nil
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
