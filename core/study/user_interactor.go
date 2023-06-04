package study

import (
	"github.com/gottenheim/ariadne/core/card"
)

type UserInteractor interface {
	ShowStudyProgress(selectedDailyCard *SelectedDailyCard, studyProgress *StudyProgress)
	AskQuestion(crd *card.Card, states []*CardState) (*CardState, error)
}
