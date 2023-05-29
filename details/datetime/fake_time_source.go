package datetime

import (
	"time"
)

type FakeTimeSource struct {
	now time.Time
}

func NewFakeTimeSource() *FakeTimeSource {
	return &FakeTimeSource{
		now: FakeNow(),
	}
}

func (s *FakeTimeSource) Today() time.Time {
	return time.Date(s.now.Year(), s.now.Month(), s.now.Day(), 0, 0, 0, 0, time.Local)
}
