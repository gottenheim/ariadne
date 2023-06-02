package datetime

import "time"

type TimeSource interface {
	Now() time.Time
	MoveNow(interval time.Duration)
}
