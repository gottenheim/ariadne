package card

import "time"

type RemindCardActivity struct {
	previousActivity CardActivity
	scheduledTo      time.Time
	executed         bool
	executionTime    time.Time
}

func CreateRemindCardActivity(previousActivity CardActivity) *RemindCardActivity {
	return &RemindCardActivity{
		previousActivity: previousActivity,
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

func (s *RemindCardActivity) IsExecuted() bool {
	return s.executed
}

func (s *RemindCardActivity) ExecutionTime() time.Time {
	return s.executionTime
}

func (s *RemindCardActivity) ScheduledTo() time.Time {
	return s.scheduledTo
}

func (s *RemindCardActivity) PreviousActivity() CardActivity {
	return s.previousActivity
}