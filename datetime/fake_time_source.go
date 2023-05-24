package datetime

import (
	"time"

	"github.com/gottenheim/ariadne/test"
)

type FakeTimeSource struct {
	now time.Time
}

func NewFakeTimeSource() *FakeTimeSource {
	return &FakeTimeSource{
		now: test.GetLocalTestTime(),
	}
}

func (s *FakeTimeSource) Today() time.Time {
	return time.Date(s.now.Year(), s.now.Month(), s.now.Day(), 0, 0, 0, 0, time.Local)
}
