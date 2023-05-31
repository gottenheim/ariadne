package card

import (
	"github.com/gottenheim/ariadne/libraries/datetime"
)

type isCardScheduledToRemindToday struct {
	timeSource datetime.TimeSource
	result     bool
}

func (s *isCardScheduledToRemindToday) OnLearnCard(learn *LearnCardActivity) error {
	return nil
}

func (s *isCardScheduledToRemindToday) OnRemindCard(remind *RemindCardActivity) error {
	s.result = !remind.IsExecuted() &&
		datetime.IsBeforeTomorrow(s.timeSource, remind.ScheduledTo())

	if s.result {
		return nil
	}

	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func IsCardScheduledToRemindToday(timeSource datetime.TimeSource, activity CardActivity) (bool, error) {
	isScheduledToRemind := &isCardScheduledToRemindToday{
		timeSource: timeSource,
		result:     false,
	}
	err := activity.Accept(isScheduledToRemind)
	if err != nil {
		return false, err
	}
	return isScheduledToRemind.result, nil
}
