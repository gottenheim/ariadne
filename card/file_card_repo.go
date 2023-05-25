package card

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gottenheim/ariadne/fs"
	"github.com/spf13/afero"
)

type FileCardRepository struct {
	fs      afero.Fs
	baseDir string
}

func NewFileCardRepository(fs afero.Fs, baseDir string) *FileCardRepository {
	return &FileCardRepository{
		fs:      fs,
		baseDir: baseDir,
	}
}

func (r *FileCardRepository) Get(relativeCardPath string) (*Card, error) {
	card, err := r.createCardFromPath(relativeCardPath)
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

	card.artifacts = artifacts

	activities, err := r.ReadCardActivities(cardPath)
	if err != nil {
		return nil, err
	}

	card.activities = activities

	return card, nil
}

func (r *FileCardRepository) Save(card *Card) error {
	err := r.assignOrderNumberIfNeeded(card)
	if err != nil {
		return err
	}

	err = r.emptyCardDirectory(card)
	if err != nil {
		return err
	}

	err = r.saveArtifactFiles(card)
	if err != nil {
		return err
	}

	return r.SaveCardActivities(card.Activities(), r.getCardPath(card))
}

func (r *FileCardRepository) assignOrderNumberIfNeeded(card *Card) error {
	if !card.HasOrderNumber() {
		cardSectionPath := r.getCardSectionPath(card)
		orderNum, err := r.getNextFreeOrderNumberInSection(cardSectionPath)

		if err != nil {
			return err
		}

		card.SetOrderNumber(orderNum)
	}

	return nil
}

func (r *FileCardRepository) getCardSectionPath(card *Card) string {
	relPath := filepath.Join(card.sections...)

	return filepath.Join(r.baseDir, relPath)
}

func (r *FileCardRepository) getNextFreeOrderNumberInSection(cardSectionPath string) (int, error) {
	maxCardNumber := 0

	err := afero.Walk(r.fs, cardSectionPath, func(filePath string, info os.FileInfo, err error) error {
		isDir, _ := afero.IsDir(r.fs, filePath)

		if isDir {
			cardDir := path.Base(filePath)
			cardNumber, err := strconv.Atoi(cardDir)
			if err != nil {
				return nil
			}

			if cardNumber > maxCardNumber {
				maxCardNumber = cardNumber
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return maxCardNumber + 1, nil
}

func (r *FileCardRepository) isCardDir(dirPath string) (bool, error) {
	answerFileExists, err := afero.Exists(r.fs, filepath.Join(dirPath, AnswerArtifactName))
	if err != nil {
		return false, err
	}

	if answerFileExists {
		return true, nil
	}

	activitiesFileExists, err := afero.Exists(r.fs, filepath.Join(dirPath, ActivitiesFileName))
	if err != nil {
		return false, err
	}

	if activitiesFileExists {
		return true, nil
	}

	dirName := path.Base(dirPath)

	orderNumber, err := strconv.Atoi(dirName)
	if err != nil {
		return false, err
	}

	return orderNumber > 0, nil
}

func (r *FileCardRepository) getCardPath(card *Card) string {
	cardSectionPath := r.getCardSectionPath(card)

	return filepath.Join(cardSectionPath, strconv.Itoa(card.orderNum))
}

func (r *FileCardRepository) emptyCardDirectory(card *Card) error {
	cardPath := r.getCardPath(card)

	cardDirExists, err := afero.Exists(r.fs, cardPath)

	if err != nil {
		return err
	}

	if cardDirExists {
		fs.EmptyDirectory(r.fs, cardPath)
	}

	return nil
}

func (r *FileCardRepository) saveArtifactFiles(card *Card) error {
	cardPath := r.getCardPath(card)

	err := r.fs.MkdirAll(cardPath, os.ModePerm)

	if err != nil {
		return err
	}

	for _, artifact := range card.artifacts {
		filePath := filepath.Join(cardPath, artifact.name)
		err = afero.WriteFile(r.fs, filePath, artifact.Content(), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *FileCardRepository) createCardFromPath(cardPath string) (*Card, error) {
	if cardPath[0] == filepath.Separator {
		cardPath = cardPath[1:]
	}
	pathItems := strings.Split(cardPath, fmt.Sprintf("%c", filepath.Separator))
	sections := pathItems[0 : len(pathItems)-1]
	orderNum, err := strconv.Atoi(pathItems[len(pathItems)-1])

	if err != nil {
		return nil, err
	}

	return &Card{
		sections: sections,
		orderNum: orderNum,
	}, nil
}

func (r *FileCardRepository) readCardArtifacts(cardPath string) ([]CardArtifact, error) {
	var artifacts []CardArtifact

	err := afero.Walk(r.fs, cardPath, func(filePath string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() && !r.isServiceFile(filePath) {
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

	return artifacts, nil
}

func (r *FileCardRepository) isServiceFile(fileName string) bool {
	return fileName[0] == '.'
}
