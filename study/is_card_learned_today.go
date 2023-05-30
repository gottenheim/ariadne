package study

import (
	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
)

type isCardLearnedToday struct {
	timeSource datetime.TimeSource
	result     bool
}

func (s *isCardLearnedToday) OnLearnCard(learn *card.LearnCardActivity) error {
	s.result = learn.IsExecuted() &&
		datetime.IsToday(s.timeSource, learn.ExecutionTime())
	return nil
}

func (s *isCardLearnedToday) OnRemindCard(remind *card.RemindCardActivity) error {
	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func IsCardLearnedToday(timeSource datetime.TimeSource, activity card.CardActivity) (bool, error) {
	isLearnedToday := &isCardLearnedToday{
		timeSource: timeSource,
		result:     false,
	}
	err := activity.Accept(isLearnedToday)
	if err != nil {
		return false, err
	}
	return isLearnedToday.result, nil
}
