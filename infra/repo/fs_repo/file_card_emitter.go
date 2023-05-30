package fs_repo

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/pipeline"
	"github.com/spf13/afero"
)

type fileCardEmitter struct {
	fs       afero.Fs
	cardsDir string
}

func NewFileCardEmitter(fs afero.Fs, cardsDir string) pipeline.Emitter[card.BriefCard] {
	return &fileCardEmitter{
		fs:       fs,
		cardsDir: cardsDir,
	}
}

func (e *fileCardEmitter) Run(ctx context.Context, output chan<- card.BriefCard) error {
	cardRepo := newFileCardRepository(e.fs, e.cardsDir)

	err := afero.Walk(e.fs, e.cardsDir, func(filePath string, info os.FileInfo, err error) error {
		isDir, _ := afero.IsDir(e.fs, filePath)

		if !isDir {
			return nil
		}

		cardDir := filePath

		isCardDir, err := e.isCardDir(cardDir)

		if err != nil {
			return err
		}

		if !isCardDir {
			return nil
		}

		cardActivities, err := cardRepo.ReadCardActivities(cardDir)

		if err != nil {
			return err
		}

		output <- card.BriefCard{
			Activities: cardActivities,
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *fileCardEmitter) isCardDir(dirPath string) (bool, error) {
	answerFileExists, err := afero.Exists(e.fs, filepath.Join(dirPath, card.AnswerArtifactName))
	if err != nil {
		return false, err
	}

	if answerFileExists {
		return true, nil
	}

	activitiesFileExists, err := afero.Exists(e.fs, filepath.Join(dirPath, ActivitiesFileName))
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
