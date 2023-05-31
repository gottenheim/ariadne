package study

import (
	"time"

	"github.com/gottenheim/ariadne/core/card"
)

type markLastActivityAsExecuted struct {
	time time.Time
}

func (s *markLastActivityAsExecuted) OnLearnCard(learn *card.LearnCardActivity) error {
	learn.MarkAsExecuted(s.time)
	return nil
}

func (s *markLastActivityAsExecuted) OnRemindCard(remind *card.RemindCardActivity) error {
	remind.MarkAsExecuted(s.time)
	return nil
}

func MarkLastActivityAsExecuted(activity card.CardActivity, time time.Time) error {
	markAsExecuted := &markLastActivityAsExecuted{
		time: time,
	}
	err := activity.Accept(markAsExecuted)
	if err != nil {
		return err
	}
	return nil
}
