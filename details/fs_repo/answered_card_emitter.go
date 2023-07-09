package fs_repo

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gottenheim/ariadne/core/card"
	"github.com/gottenheim/ariadne/libraries/pipeline"
	"github.com/spf13/afero"
)

type fileCardEmitter struct {
	fs       afero.Fs
	cardRepo *fileCardRepository
	cardsDir string
}

func NewAnsweredCardEmitter(fs afero.Fs, cardRepo *fileCardRepository, cardsDir string) pipeline.Emitter[card.BriefCard] {
	return &fileCardEmitter{
		fs:       fs,
		cardRepo: cardRepo,
		cardsDir: cardsDir,
	}
}

func (e *fileCardEmitter) Run(ctx context.Context, output chan<- card.BriefCard) error {
	err := afero.Walk(e.fs, e.cardsDir, func(filePath string, info os.FileInfo, err error) error {
		isDir, _ := afero.IsDir(e.fs, filePath)

		if !isDir {
			return nil
		}

		cardDir := filePath

		containsAnswerFile, err := e.containsAnswerFile(cardDir)

		if err != nil {
			return err
		}

		if !containsAnswerFile {
			return nil
		}

		cardActivities, err := e.cardRepo.ReadCardActivities(cardDir)

		if err != nil {
			return fmt.Errorf("failed to read activities for card %v", cardDir)
		}

		section, entry := e.cardRepo.GetCardPathSection(cardDir), e.cardRepo.GetCardPathEntry(cardDir)

		briefCard := card.BriefCard{
			Section:    section,
			Entry:      entry,
			Activities: cardActivities,
		}

		if !pipeline.WriteToChannel[card.BriefCard](ctx, output, briefCard) {
			return io.EOF
		}

		return nil
	})

	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (e *fileCardEmitter) containsAnswerFile(dirPath string) (bool, error) {
	answerFileExists, err := afero.Exists(e.fs, filepath.Join(dirPath, card.AnswerArtifactName))
	if err != nil {
		return false, err
	}

	return answerFileExists, nil
}
