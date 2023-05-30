package card

import (
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type isCardLearnedToday struct {
	timeSource datetime.TimeSource
	result     bool
}

func (s *isCardLearnedToday) OnLearnCard(learn *LearnCardActivity) error {
	s.result = learn.IsExecuted() &&
		datetime.IsToday(s.timeSource, learn.ExecutionTime())
	return nil
}

func (s *isCardLearnedToday) OnRemindCard(remind *RemindCardActivity) error {
	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func IsCardLearnedToday(timeSource datetime.TimeSource, activity CardActivity) (bool, error) {
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
