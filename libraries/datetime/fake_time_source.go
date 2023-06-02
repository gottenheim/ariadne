package datetime

import (
	"time"
)

type FakeTimeSource struct {
	now time.Time
}

func NewFakeTimeSource() TimeSource {
	return &FakeTimeSource{
		now: FakeNow(),
	}
}

func (s *FakeTimeSource) Now() time.Time {
	return s.now
}

func (s *FakeTimeSource) MoveNow(d time.Duration) {
	s.now = s.now.Add(d)
}
