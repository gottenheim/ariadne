package card

import "github.com/gottenheim/ariadne/details/datetime"

type isCardLearnedToday struct {
	timeSource datetime.TimeSource
	result     bool
}

func (s *isCardLearnedToday) OnLearnCard(learn *LearnCardActivity) error {
	s.result = learn.executed &&
		datetime.IsToday(s.timeSource, learn.executionTime)
	return nil
}

func (s *isCardLearnedToday) OnRemindCard(remind *RemindCardActivity) error {
	if remind.previousActivity == nil {
		return nil
	}

	return remind.previousActivity.Accept(s)
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
