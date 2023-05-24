package card

import "time"

type RemindCardActivity struct {
	previousActivity CardActivity
	scheduledTo      time.Time
	executed         bool
	executionTime    time.Time
}

func CreateRemindCardActivity(previousState CardActivity) *RemindCardActivity {
	return &RemindCardActivity{
		previousActivity: previousState,
	}
}

func (s *RemindCardActivity) Accept(visitor CardActivityVisitor) error {
	return visitor.OnRemindCard(s)
}

func (s *RemindCardActivity) ScheduleTo(scheduledTo time.Time) {
	s.scheduledTo = scheduledTo
}

func (s *RemindCardActivity) MarkAsExecuted(executionTime time.Time) {
	s.executed = true
	s.executionTime = executionTime
}
