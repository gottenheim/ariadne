package study

import (
	"time"

	"github.com/gottenheim/ariadne/core/card"
)

const initialEasinessFactor = 2.5

type CardRepetitionParams struct {
	EasinessFactor   float64
	RepetitionNumber int
	Interval         time.Duration
}

func (s *CardRepetitionParams) OnLearnCard(learn *card.LearnCardActivity) error {
	return nil
}

func (s *CardRepetitionParams) OnRemindCard(remind *card.RemindCardActivity) error {
	s.EasinessFactor = remind.EasinessFactor()
	s.RepetitionNumber = remind.RepetitionNumber()
	s.Interval = remind.Interval()
	return nil
}

func GetCardRepetitionParams(activity card.CardActivity) (*CardRepetitionParams, error) {
	repetitionParams := &CardRepetitionParams{
		EasinessFactor: initialEasinessFactor,
	}
	err := activity.Accept(repetitionParams)
	if err != nil {
		return nil, err
	}
	return repetitionParams, nil
}
