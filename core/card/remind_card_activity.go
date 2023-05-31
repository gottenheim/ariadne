package card

import "time"

type RemindCardActivity struct {
	previousActivity CardActivity
	scheduledTo      time.Time
	executed         bool
	executionTime    time.Time
	easinessFactor   float64
	repetitionNumber int
	interval         int
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

func (s *RemindCardActivity) EasinessFactor() float64 {
	return s.easinessFactor
}

func (s *RemindCardActivity) SetEasinessFactor(easinessFactor float64) {
	s.easinessFactor = easinessFactor
}

func (s *RemindCardActivity) RepetitionNumber() int {
	return s.repetitionNumber
}

func (s *RemindCardActivity) SetRepetitionNumber(repetitionNumber int) {
	s.repetitionNumber = repetitionNumber
}

func (s *RemindCardActivity) Interval() int {
	return s.interval
}

func (s *RemindCardActivity) SetInterval(interval int) {
	s.interval = interval
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
