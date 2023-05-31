package datetime

import "time"

type TimeSource interface {
	Today() time.Time
	Now() time.Time
}
