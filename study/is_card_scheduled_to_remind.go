package study

import (
	"github.com/gottenheim/ariadne/card"
	"github.com/gottenheim/ariadne/details/datetime"
)

type isCardScheduledToRemindToday struct {
	timeSource datetime.TimeSource
	result     bool
}

func (s *isCardScheduledToRemindToday) OnLearnCard(learn *card.LearnCardActivity) error {
	return nil
}

func (s *isCardScheduledToRemindToday) OnRemindCard(remind *card.RemindCardActivity) error {
	s.result = !remind.IsExecuted() &&
		datetime.IsBeforeTomorrow(s.timeSource, remind.ScheduledTo())

	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func IsCardScheduledToRemindToday(timeSource datetime.TimeSource, activity card.CardActivity) (bool, error) {
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
