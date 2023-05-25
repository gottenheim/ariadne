package card

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const ActivitiesFileName = ".activities"

type FileCardActivityRepository struct {
	fs afero.Fs
}

func NewFileCardActivityRepository(fs afero.Fs) *FileCardActivityRepository {
	return &FileCardActivityRepository{
		fs: fs,
	}
}

func (r *FileCardActivityRepository) ReadCardActivities(cardPath string) (CardActivity, error) {
	activitiesFilePath := r.getActivitiesFilePath(cardPath)

	fileExists, activitiesBinary, err := r.readActivitiesFromFile(activitiesFilePath)
	if err != nil {
		return nil, err
	}

	if !fileExists {
		return CreateLearnCardActivity(), nil
	}

	return DeserializeCardActivityChain(activitiesBinary)
}

func (r *FileCardActivityRepository) SaveCardActivities(cardActivity CardActivity, cardPath string) error {
	activitiesFilePath := r.getActivitiesFilePath(cardPath)

	r.removeActivitiesFileIfExists(activitiesFilePath)

	isNewCard, err := IsNewCard(cardActivity)

	if err != nil {
		return err
	}

	if isNewCard {
		return nil
	}

	activitiesBinary, err := SerializeCardActivityChain(cardActivity)
	if err != nil {
		return err
	}

	return r.writeActivitiesToFile(activitiesFilePath, activitiesBinary)
}

func (r *FileCardActivityRepository) getActivitiesFilePath(cardPath string) string {
	return filepath.Join(cardPath, ActivitiesFileName)
}

func (r *FileCardActivityRepository) readActivitiesFromFile(activitiesFilePath string) (bool, []byte, error) {
	exists, err := afero.Exists(r.fs, activitiesFilePath)
	if err != nil {
		return false, nil, err
	}

	if !exists {
		return false, nil, nil
	}

	activitiesBinary, err := afero.ReadFile(r.fs, activitiesFilePath)
	if err != nil {
		return false, nil, err
	}

	return true, activitiesBinary, nil
}

func (r *FileCardActivityRepository) removeActivitiesFileIfExists(activitiesFilePath string) error {
	exists, err := afero.Exists(r.fs, activitiesFilePath)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	err = r.fs.Remove(activitiesFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileCardActivityRepository) writeActivitiesToFile(activitiesFilePath string, progress []byte) error {
	return afero.WriteFile(r.fs, activitiesFilePath, progress, os.ModePerm)
}
