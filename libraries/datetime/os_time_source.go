package datetime

import (
	"time"
)

type OsTimeSource struct {
}

func NewOsTimeSource() TimeSource {
	return &OsTimeSource{}
}

func (s *OsTimeSource) Now() time.Time {
	return time.Now()
}
