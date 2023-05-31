package study

import (
	"github.com/gottenheim/ariadne/core/card"
)

const initialEasinessFactor = 2.5

type cardRepetitionParams struct {
	easinessFactor   float64
	repetitionNumber int
	interval         int
}

func (s *cardRepetitionParams) OnLearnCard(learn *card.LearnCardActivity) error {
	return nil
}

func (s *cardRepetitionParams) OnRemindCard(remind *card.RemindCardActivity) error {
	s.easinessFactor = remind.EasinessFactor()
	s.repetitionNumber = remind.RepetitionNumber()
	s.interval = remind.Interval()
	return nil
}

func GetCardRepetitionParams(activity card.CardActivity) (*cardRepetitionParams, error) {
	repetitionParams := &cardRepetitionParams{
		easinessFactor: initialEasinessFactor,
	}
	err := activity.Accept(repetitionParams)
	if err != nil {
		return nil, err
	}
	return repetitionParams, nil
}
