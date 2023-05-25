package card_repo

import (
	"os"
	"path/filepath"

	"github.com/gottenheim/ariadne/card"
	"github.com/spf13/afero"
)

const ActivitiesFileName = ".activities"

func (r *FileCardRepository) ReadCardActivities(cardPath string) (card.CardActivity, error) {
	activitiesFilePath := r.getActivitiesFilePath(cardPath)

	fileExists, activitiesBinary, err := r.readActivitiesFromFile(activitiesFilePath)
	if err != nil {
		return nil, err
	}

	if !fileExists {
		return card.CreateLearnCardActivity(), nil
	}

	return card.DeserializeCardActivityChain(activitiesBinary)
}

func (r *FileCardRepository) SaveCardActivities(cardActivity card.CardActivity, cardPath string) error {
	activitiesFilePath := r.getActivitiesFilePath(cardPath)

	r.removeActivitiesFileIfExists(activitiesFilePath)

	isNewCard, err := card.IsNewCard(cardActivity)

	if err != nil {
		return err
	}

	if isNewCard {
		return nil
	}

	activitiesBinary, err := card.SerializeCardActivityChain(cardActivity)
	if err != nil {
		return err
	}

	return r.writeActivitiesToFile(activitiesFilePath, activitiesBinary)
}

func (r *FileCardRepository) getActivitiesFilePath(cardPath string) string {
	return filepath.Join(cardPath, ActivitiesFileName)
}

func (r *FileCardRepository) readActivitiesFromFile(activitiesFilePath string) (bool, []byte, error) {
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

func (r *FileCardRepository) removeActivitiesFileIfExists(activitiesFilePath string) error {
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

func (r *FileCardRepository) writeActivitiesToFile(activitiesFilePath string, progress []byte) error {
	return afero.WriteFile(r.fs, activitiesFilePath, progress, os.ModePerm)
}
