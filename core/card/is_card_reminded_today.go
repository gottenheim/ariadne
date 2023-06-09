package card

import (
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type isCardRemindedToday struct {
	timeSource datetime.TimeSource
	result     bool
}

func (s *isCardRemindedToday) OnLearnCard(learn *LearnCardActivity) error {
	return nil
}

func (s *isCardRemindedToday) OnRemindCard(remind *RemindCardActivity) error {
	s.result = remind.IsExecuted() &&
		datetime.IsToday(s.timeSource, remind.ExecutionTime())

	if s.result {
		return nil
	}

	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func IsCardRemindedToday(timeSource datetime.TimeSource, activity CardActivity) (bool, error) {
	isRemindedToday := &isCardRemindedToday{
		timeSource: timeSource,
		result:     false,
	}
	err := activity.Accept(isRemindedToday)
	if err != nil {
		return false, err
	}
	return isRemindedToday.result, nil
}
