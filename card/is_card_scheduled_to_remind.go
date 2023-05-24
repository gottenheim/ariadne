package card

import "github.com/gottenheim/ariadne/datetime"

type isCardScheduledToRemindToday struct {
	timeSource datetime.TimeSource
	result     bool
}

func (s *isCardScheduledToRemindToday) OnLearnCard(learn *LearnCardActivity) error {
	return nil
}

func (s *isCardScheduledToRemindToday) OnRemindCard(remind *RemindCardActivity) error {
	s.result = !remind.executed &&
		datetime.IsBeforeTomorrow(s.timeSource, remind.scheduledTo)

	if remind.previousActivity == nil {
		return nil
	}

	return remind.previousActivity.Accept(s)
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
