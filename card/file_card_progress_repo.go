package card

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/gottenheim/ariadne/config"
	"github.com/spf13/afero"
)

const ProgressFileName = ".progress"

type FileCardProgressRepository struct {
	fs afero.Fs
}

func NewFileCardProgressRepository(fs afero.Fs) *FileCardProgressRepository {
	return &FileCardProgressRepository{
		fs: fs,
	}
}

func (r *FileCardProgressRepository) ReadCardProgress(cardPath string) (*CardProgress, error) {
	progressFilePath := r.getProgressFilePath(cardPath)

	fileExists, progressBinary, err := r.readProgressFromFile(progressFilePath)
	if err != nil {
		return nil, err
	}

	if !fileExists {
		return r.getNewCardProgress(), nil
	}

	return r.deserializeProgress(progressBinary)
}

func (r *FileCardProgressRepository) SaveCardProgress(cardProgress *CardProgress, cardPath string) error {
	progressFilePath := r.getProgressFilePath(cardPath)

	r.removeProgressFileIfExists(progressFilePath)

	if cardProgress.IsNew() {
		return nil
	}

	progressBinary, err := r.serializeProgress(cardProgress)
	if err != nil {
		return err
	}

	return r.writeProgressToFile(progressFilePath, progressBinary)
}

func (r *FileCardProgressRepository) getProgressFilePath(cardPath string) string {
	return filepath.Join(cardPath, ProgressFileName)
}

func (r *FileCardProgressRepository) readProgressFromFile(progressFilePath string) (bool, []byte, error) {
	exists, err := afero.Exists(r.fs, progressFilePath)
	if err != nil {
		return false, nil, err
	}

	if !exists {
		return false, nil, nil
	}

	progressBinary, err := afero.ReadFile(r.fs, progressFilePath)
	if err != nil {
		return false, nil, err
	}

	return true, progressBinary, nil
}

func (r *FileCardProgressRepository) deserializeProgress(progressBinary []byte) (*CardProgress, error) {
	cfg, err := config.FromYamlReader(bytes.NewReader(progressBinary))

	if err != nil {
		return nil, err
	}

	progressModel := &CardProgressWriteModel{}
	err = cfg.Materialize(progressModel)

	if err != nil {
		return nil, err
	}

	return progressModel.ToCardProgress(), nil
}

func (r *FileCardProgressRepository) removeProgressFileIfExists(progressFilePath string) error {
	exists, err := afero.Exists(r.fs, progressFilePath)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	err = r.fs.Remove(progressFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileCardProgressRepository) serializeProgress(cardProgress *CardProgress) ([]byte, error) {
	cardProgressModel := cardProgress.ToWriteModel()

	return config.SerializeToYaml(cardProgressModel)
}

func (r *FileCardProgressRepository) writeProgressToFile(progressFilePath string, progress []byte) error {
	return afero.WriteFile(r.fs, progressFilePath, progress, os.ModePerm)
}

func (r *FileCardProgressRepository) getNewCardProgress() *CardProgress {
	return GetNewCardProgress()
}
