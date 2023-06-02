package datetime

import "time"

type TimeSource interface {
	Now() time.Time
}
