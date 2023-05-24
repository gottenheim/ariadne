package card

import "time"

type LearnCardActivity struct {
	executed      bool
	executionTime time.Time
}

func CreateLearnCardActivity() *LearnCardActivity {
	return &LearnCardActivity{}
}

func (s *LearnCardActivity) Accept(visitor CardActivityVisitor) error {
	return visitor.OnLearnCard(s)
}

func (s *LearnCardActivity) MarkAsExecuted(executionTime time.Time) {
	s.executed = true
	s.executionTime = executionTime
}
