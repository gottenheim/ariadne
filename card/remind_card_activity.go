package card

import "time"

type RemindCardActivity struct {
	previousActivity CardActivity
	scheduledTo      time.Time
	executed         bool
	executionTime    time.Time
}

func CreateRemindCardActivity(scheduledTo time.Time, previousState CardActivity) CardActivity {
	return &RemindCardActivity{
		scheduledTo:      scheduledTo,
		previousActivity: previousState,
	}
}

func (s *RemindCardActivity) Accept(visitor CardActivityVisitor) error {
	return visitor.OnRemindCard(s)
}
