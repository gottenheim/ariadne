package card

import (
	"time"

	"github.com/gottenheim/ariadne/libraries/datetime"
)

type cardTodayReminderTime struct {
	timeSource datetime.TimeSource
	result     time.Time
}

func (s *cardTodayReminderTime) OnLearnCard(learn *LearnCardActivity) error {
	return nil
}

func (s *cardTodayReminderTime) OnRemindCard(remind *RemindCardActivity) error {
	if !remind.IsExecuted() &&
		datetime.IsToday(s.timeSource, remind.ScheduledTo()) {
		s.result = remind.ScheduledTo()
		return nil
	}

	if remind.PreviousActivity() == nil {
		return nil
	}

	return remind.PreviousActivity().Accept(s)
}

func GetTimeToRemindToday(timeSource datetime.TimeSource, activity CardActivity) (time.Time, error) {
	reminderTime := &cardTodayReminderTime{
		timeSource: timeSource,
	}
	err := activity.Accept(reminderTime)
	if err != nil {
		return time.Time{}, err
	}
	return reminderTime.result, nil
}
