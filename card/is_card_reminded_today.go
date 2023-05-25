package card

import "github.com/gottenheim/ariadne/details/datetime"

type isCardRemindedToday struct {
	timeSource datetime.TimeSource
	result     bool
}

func (s *isCardRemindedToday) OnLearnCard(learn *LearnCardActivity) error {
	return nil
}

func (s *isCardRemindedToday) OnRemindCard(remind *RemindCardActivity) error {
	s.result = remind.executed &&
		datetime.IsToday(s.timeSource, remind.executionTime)

	if s.result {
		return nil
	}

	if remind.previousActivity == nil {
		return nil
	}

	return remind.previousActivity.Accept(s)
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
